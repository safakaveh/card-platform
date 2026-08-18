package getdata

import (
	"context"
	"database/sql"
	"fmt"
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
