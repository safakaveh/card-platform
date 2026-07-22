package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

func Open(ctx context.Context, databasePath string) (*sql.DB, error) {
	if databasePath == "" {
		return nil, fmt.Errorf("database path is empty")
	}

	absolutePath, err := filepath.Abs(databasePath)
	if err != nil {
		return nil, fmt.Errorf("resolve database path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	dsn := sqliteDSN(absolutePath)

	database, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	// SQLite معمولاً با تعداد connection محدود پایدارتر است.
	database.SetMaxOpenConns(4)
	database.SetMaxIdleConns(4)
	database.SetConnMaxLifetime(0)
	database.SetConnMaxIdleTime(5 * time.Minute)

	if err := database.PingContext(ctx); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	if err := Migrate(ctx, database); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("migrate sqlite database: %w", err)
	}

	return database, nil
}

func sqliteDSN(databasePath string) string {
	uri := &url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(databasePath),
	}

	query := uri.Query()

	// برای تمام connectionهای pool اعمال می‌شود.
	query.Add("_pragma", "foreign_keys(1)")
	query.Add("_pragma", "busy_timeout(5000)")
	query.Add("_pragma", "journal_mode(WAL)")
	query.Add("_pragma", "synchronous(NORMAL)")

	uri.RawQuery = query.Encode()

	return uri.String()
}
