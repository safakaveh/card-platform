package getdata

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type Service struct{ database *sql.DB }

func NewService(database *sql.DB) *Service { return &Service{database: database} }

func (s *Service) PendingMifare(ctx context.Context, limit int) (PendingResponse, error) {
	if limit < 1 || limit > 1000 {
		limit = 100
	}
	const query = `
		SELECT m.uuid_card, c.uuid_order, o.order_name, m.block_no, m.created_at
		FROM mifare_data m
		INNER JOIN cards c ON c.uuid = m.uuid_card
		INNER JOIN orders o ON o.uuid = c.uuid_order
		WHERE m.content = ?
		ORDER BY m.created_at ASC
		LIMIT ?`
	rows, err := s.database.QueryContext(ctx, query, []byte("UUID_PENDING"), limit)
	if err != nil {
		return PendingResponse{}, fmt.Errorf("query pending mifare data: %w", err)
	}
	defer rows.Close()
	items := make([]PendingCard, 0)
	for rows.Next() {
		var item PendingCard
		if err := rows.Scan(&item.CardUUID, &item.OrderUUID, &item.OrderName, &item.BlockNo, &item.CreatedAt); err != nil {
			return PendingResponse{}, fmt.Errorf("scan pending mifare data: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return PendingResponse{}, fmt.Errorf("iterate pending mifare data: %w", err)
	}
	return PendingResponse{Count: len(items), Items: items}, nil
}

func (s *Service) ReadNext(ctx context.Context, orderName string, limit int) (ReadResponse, error) {
	orderName = strings.TrimSpace(orderName)
	if orderName == "" {
		return ReadResponse{}, ErrOrderNameRequired
	}
	if limit < 1 || limit > 100 {
		limit = 1
	}

	tx, err := s.database.BeginTx(ctx, nil)
	if err != nil {
		return ReadResponse{}, fmt.Errorf("begin read transaction: %w", err)
	}
	defer tx.Rollback()

	var orderID string
	if err := tx.QueryRowContext(ctx, `SELECT uuid FROM orders WHERE order_name = ?`, orderName).Scan(&orderID); err != nil {
		if err == sql.ErrNoRows {
			return ReadResponse{}, ErrOrderNotFound
		}
		return ReadResponse{}, fmt.Errorf("find order: %w", err)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT c.uuid, c.uuid_order, o.order_name, c.created_at
		FROM cards c
		INNER JOIN orders o ON o.uuid = c.uuid_order
		WHERE c.uuid_order = ? AND c.read_at IS NULL
		ORDER BY c.created_at ASC, c.uuid ASC
		LIMIT ?`, orderID, limit)
	if err != nil {
		return ReadResponse{}, fmt.Errorf("select unread cards: %w", err)
	}
	items, err := scanReadItems(ctx, tx, rows, true)
	rows.Close()
	if err != nil {
		return ReadResponse{}, err
	}
	now := currentMillis()
	for i := range items {
		if _, err := tx.ExecContext(ctx, `UPDATE cards SET read_at = ?, updated_at = ? WHERE uuid = ? AND read_at IS NULL`, now, now, items[i].CardUUID); err != nil {
			return ReadResponse{}, fmt.Errorf("mark card read: %w", err)
		}
		items[i].ReadAt = &now
	}
	if err := tx.Commit(); err != nil {
		return ReadResponse{}, fmt.Errorf("commit read transaction: %w", err)
	}
	return ReadResponse{Count: len(items), Items: items}, nil
}

func (s *Service) ReadReport(ctx context.Context, orderName, status string, limit int) (ReadReportResponse, error) {
	if limit < 1 || limit > 1000 {
		limit = 100
	}
	query := `
		SELECT c.uuid, c.uuid_order, o.order_name, c.created_at, c.read_at
		FROM cards c INNER JOIN orders o ON o.uuid = c.uuid_order`
	args := make([]any, 0, 2)
	conditions := make([]string, 0, 2)
	if strings.TrimSpace(orderName) != "" {
		conditions = append(conditions, "o.order_name = ?")
		args = append(args, strings.TrimSpace(orderName))
	}
	switch status {
	case "read":
		conditions = append(conditions, "c.read_at IS NOT NULL")
	case "unread":
		conditions = append(conditions, "c.read_at IS NULL")
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY c.created_at ASC, c.uuid ASC LIMIT ?"
	args = append(args, limit)
	rows, err := s.database.QueryContext(ctx, query, args...)
	if err != nil {
		return ReadReportResponse{}, fmt.Errorf("query read report: %w", err)
	}
	defer rows.Close()
	items := make([]ReadReportItem, 0)
	for rows.Next() {
		var item ReadReportItem
		var readAt sql.NullInt64
		if err := rows.Scan(&item.CardUUID, &item.OrderUUID, &item.OrderName, &item.CreatedAt, &readAt); err != nil {
			return ReadReportResponse{}, fmt.Errorf("scan read report: %w", err)
		}
		if readAt.Valid {
			item.ReadAt = &readAt.Int64
			item.Read = true
		}
		if err := loadCardData(ctx, s.database, &item.ReadItem); err != nil {
			return ReadReportResponse{}, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return ReadReportResponse{}, fmt.Errorf("iterate read report: %w", err)
	}
	return ReadReportResponse{Count: len(items), Items: items}, nil
}

func (s *Service) ResetRead(ctx context.Context, cardID string) error {
	result, err := s.database.ExecContext(ctx, `UPDATE cards SET read_at = NULL, updated_at = ? WHERE uuid = ?`, currentMillis(), cardID)
	if err != nil {
		return fmt.Errorf("reset card read status: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read reset result: %w", err)
	}
	if n == 0 {
		return ErrCardNotFound
	}
	return nil
}

func scanReadItems(ctx context.Context, tx *sql.Tx, rows *sql.Rows, loadData bool) ([]ReadItem, error) {
	items := make([]ReadItem, 0)
	for rows.Next() {
		var item ReadItem
		if err := rows.Scan(&item.CardUUID, &item.OrderUUID, &item.OrderName, &item.Sequence); err != nil {
			return nil, fmt.Errorf("scan unread cards: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate unread cards: %w", err)
	}
	if loadData {
		for i := range items {
			if err := loadCardData(ctx, tx, &items[i]); err != nil {
				return nil, err
			}
		}
	}
	return items, nil
}

type queryer interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
}

func loadCardData(ctx context.Context, q queryer, item *ReadItem) error {
	laserRows, err := q.QueryContext(ctx, `SELECT side, row_no, content_type, content FROM laser_data WHERE uuid_card = ? ORDER BY side, row_no`, item.CardUUID)
	if err != nil {
		return fmt.Errorf("query laser data: %w", err)
	}
	defer laserRows.Close()
	for laserRows.Next() {
		var v LaserValue
		var content []byte
		if err := laserRows.Scan(&v.Side, &v.Row, &v.ContentType, &content); err != nil {
			return fmt.Errorf("scan laser data: %w", err)
		}
		v.Content = string(content)
		item.Laser = append(item.Laser, v)
	}
	magnetRows, err := q.QueryContext(ctx, `SELECT track_no, content FROM magnet_data WHERE uuid_card = ? ORDER BY track_no`, item.CardUUID)
	if err != nil {
		return fmt.Errorf("query magnet data: %w", err)
	}
	defer magnetRows.Close()
	for magnetRows.Next() {
		var v MagnetValue
		var content []byte
		if err := magnetRows.Scan(&v.Track, &content); err != nil {
			return fmt.Errorf("scan magnet data: %w", err)
		}
		v.Content = string(content)
		item.Magnet = append(item.Magnet, v)
	}
	mifareRows, err := q.QueryContext(ctx, `SELECT block_no, content FROM mifare_data WHERE uuid_card = ? ORDER BY block_no`, item.CardUUID)
	if err != nil {
		return fmt.Errorf("query mifare data: %w", err)
	}
	defer mifareRows.Close()
	for mifareRows.Next() {
		var v MifareValue
		var content []byte
		if err := mifareRows.Scan(&v.Block, &content); err != nil {
			return fmt.Errorf("scan mifare data: %w", err)
		}
		v.Content = string(content)
		item.Mifare = append(item.Mifare, v)
	}
	return nil
}

func currentMillis() int64 {
	return time.Now().UnixMilli()
}
