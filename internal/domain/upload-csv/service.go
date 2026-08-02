package uploadcsv

import (
	"context"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	database *sql.DB
}

func NewService(database *sql.DB) *Service {
	return &Service{database: database}
}

func (s *Service) Import(
	ctx context.Context,
	orderName string,
	fileName string,
	reader io.Reader,
) (ImportResult, error) {
	orderName = strings.TrimSpace(orderName)
	if orderName == "" {
		orderName = strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
	}
	if orderName == "" {
		return ImportResult{}, ErrEmptyOrderName
	}

	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = 0
	csvReader.ReuseRecord = true
	csvReader.TrimLeadingSpace = true

	headers, err := csvReader.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return ImportResult{}, fmt.Errorf("%w: file is empty", ErrInvalidCSV)
		}
		return ImportResult{}, fmt.Errorf("%w: read header: %v", ErrInvalidCSV, err)
	}

	mappings, frontColumns, backColumns, hasUID, err := mapHeaders(headers)
	if err != nil {
		return ImportResult{}, err
	}

	transaction, err := s.database.BeginTx(ctx, nil)
	if err != nil {
		return ImportResult{}, fmt.Errorf("begin import transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			_ = transaction.Rollback()
		}
	}()

	orderID := uuid.NewString()
	now := time.Now().UnixMilli()
	description := "CSV:" + filepath.Base(fileName)

	_, err = transaction.ExecContext(
		ctx,
		`INSERT INTO orders (
			uuid, order_name, status, description, order_date, created_at, updated_at
		) VALUES (?, ?, 'pending', ?, ?, ?, ?)`,
		orderID,
		orderName,
		description,
		now,
		now,
		now,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			return ImportResult{}, ErrOrderNameConflict
		}
		return ImportResult{}, fmt.Errorf("create import order: %w", err)
	}

	cardStatement, err := transaction.PrepareContext(ctx, `
		INSERT INTO cards (
			uuid, uuid_order, has_laser, has_magnet, has_mifare_uid,
			has_mifare, has_java_card, has_press, has_temperature,
			is_done, description, created_at, updated_at
		) VALUES (?, ?, 1, 0, ?, 0, 0, 0, 0, 0, ?, ?, ?)
	`)
	if err != nil {
		return ImportResult{}, fmt.Errorf("prepare card insert: %w", err)
	}
	defer cardStatement.Close()

	laserStatement, err := transaction.PrepareContext(ctx, `
		INSERT INTO laser_data (
			uuid, uuid_card, side, row_no, content_type, content, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return ImportResult{}, fmt.Errorf("prepare laser data insert: %w", err)
	}
	defer laserStatement.Close()

	var rowsImported int64
	for csvRowNumber := int64(2); ; csvRowNumber++ {
		record, readErr := csvReader.Read()
		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			return ImportResult{}, fmt.Errorf(
				"%w: row %d: %v",
				ErrInvalidCSV,
				csvRowNumber,
				readErr,
			)
		}

		cardID := uuid.NewString()
		cardDescription := fmt.Sprintf("CSV row %d", csvRowNumber)
		hasUIDValue := 0
		if hasUID {
			hasUIDValue = 1
		}

		if _, err := cardStatement.ExecContext(
			ctx,
			cardID,
			orderID,
			hasUIDValue,
			cardDescription,
			now,
			now,
		); err != nil {
			return ImportResult{}, fmt.Errorf("insert card at row %d: %w", csvRowNumber, err)
		}

		for _, mapping := range mappings {
			if _, err := laserStatement.ExecContext(
				ctx,
				uuid.NewString(),
				cardID,
				mapping.side,
				mapping.rowNumber,
				mapping.contentType,
				[]byte(record[mapping.index]),
				now,
			); err != nil {
				return ImportResult{}, fmt.Errorf(
					"insert %s at row %d: %w",
					mapping.header,
					csvRowNumber,
					err,
				)
			}
		}

		rowsImported++
	}

	if rowsImported == 0 {
		return ImportResult{}, fmt.Errorf("%w: file has no data rows", ErrInvalidCSV)
	}

	if err := transaction.Commit(); err != nil {
		return ImportResult{}, fmt.Errorf("commit import: %w", err)
	}
	committed = true

	return ImportResult{
		UUID:         orderID,
		OrderName:    orderName,
		FileName:     filepath.Base(fileName),
		RowsImported: rowsImported,
		FrontColumns: frontColumns,
		BackColumns:  backColumns,
		HasUID:       hasUID,
	}, nil
}

func (s *Service) List(ctx context.Context, limit int) ([]ImportSummary, error) {
	if limit < 1 || limit > 200 {
		limit = 50
	}

	rows, err := s.database.QueryContext(ctx, `
		SELECT
			o.uuid,
			o.order_name,
			o.status,
			CASE
				WHEN o.description LIKE 'CSV:%' THEN substr(o.description, 5)
				ELSE ''
			END AS file_name,
			COUNT(c.uuid) AS card_count,
			o.created_at
		FROM orders o
		LEFT JOIN cards c ON c.uuid_order = o.uuid
		GROUP BY o.uuid
		ORDER BY o.created_at DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("list imports: %w", err)
	}
	defer rows.Close()

	imports := make([]ImportSummary, 0)
	for rows.Next() {
		var item ImportSummary
		if err := rows.Scan(
			&item.UUID,
			&item.OrderName,
			&item.Status,
			&item.FileName,
			&item.CardCount,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan import: %w", err)
		}
		imports = append(imports, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate imports: %w", err)
	}
	return imports, nil
}

func (s *Service) Get(ctx context.Context, importID string) (ImportDetails, error) {
	var details ImportDetails
	err := s.database.QueryRowContext(ctx, `
		SELECT
			o.uuid,
			o.order_name,
			o.status,
			CASE
				WHEN o.description LIKE 'CSV:%' THEN substr(o.description, 5)
				ELSE ''
			END AS file_name,
			COUNT(DISTINCT c.uuid) AS card_count,
			o.created_at,
			COUNT(DISTINCT CASE WHEN l.side = 'front' THEN l.row_no END),
			COUNT(DISTINCT CASE WHEN l.side = 'back' THEN l.row_no END),
			COUNT(CASE WHEN l.content_type = 'uid' THEN 1 END)
		FROM orders o
		LEFT JOIN cards c ON c.uuid_order = o.uuid
		LEFT JOIN laser_data l ON l.uuid_card = c.uuid
		WHERE o.uuid = ?
		GROUP BY o.uuid
	`, importID).Scan(
		&details.UUID,
		&details.OrderName,
		&details.Status,
		&details.FileName,
		&details.CardCount,
		&details.CreatedAt,
		&details.FrontColumns,
		&details.BackColumns,
		&details.UIDFields,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ImportDetails{}, sql.ErrNoRows
		}
		return ImportDetails{}, fmt.Errorf("get import: %w", err)
	}
	return details, nil
}

func (s *Service) Delete(ctx context.Context, importID string) error {
	result, err := s.database.ExecContext(ctx, `DELETE FROM orders WHERE uuid = ?`, importID)
	if err != nil {
		return fmt.Errorf("delete import: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read deleted rows: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func mapHeaders(
	headers []string,
) ([]columnMapping, []string, []string, bool, error) {
	mappings := make([]columnMapping, 0)
	frontColumns := make([]string, 0)
	backColumns := make([]string, 0)
	seen := make(map[string]struct{})
	frontRow := int64(0)
	backRow := int64(0)
	hasUID := false

	for index, rawHeader := range headers {
		header := strings.TrimSpace(strings.TrimPrefix(rawHeader, "\ufeff"))
		normalized := strings.ToLower(header)

		side := ""
		contentType := ""
		rowNumber := int64(0)
		switch {
		case strings.HasPrefix(normalized, "frn_"):
			side = "front"
			contentType = strings.TrimSpace(header[len("frn_"):])
			frontRow++
			rowNumber = frontRow
		case strings.HasPrefix(normalized, "bck_"):
			side = "back"
			contentType = strings.TrimSpace(header[len("bck_"):])
			backRow++
			rowNumber = backRow
		default:
			continue
		}

		if contentType == "" {
			return nil, nil, nil, false, fmt.Errorf(
				"%w: column %q has no name after its prefix",
				ErrInvalidCSV,
				header,
			)
		}
		if _, exists := seen[normalized]; exists {
			return nil, nil, nil, false, fmt.Errorf(
				"%w: %s",
				ErrDuplicateColumn,
				header,
			)
		}
		seen[normalized] = struct{}{}

		isUID := strings.EqualFold(contentType, "uid")
		if isUID {
			contentType = "uid"
			hasUID = true
		}

		mappings = append(mappings, columnMapping{
			index:       index,
			side:        side,
			rowNumber:   rowNumber,
			contentType: contentType,
			header:      header,
			isUID:       isUID,
		})
		if side == "front" {
			frontColumns = append(frontColumns, header)
		} else {
			backColumns = append(backColumns, header)
		}
	}

	if len(mappings) == 0 {
		return nil, nil, nil, false, ErrNoMappedColumns
	}
	return mappings, frontColumns, backColumns, hasUID, nil
}
