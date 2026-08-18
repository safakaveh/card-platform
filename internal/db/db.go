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

	if err := ensureWritable(absolutePath); err != nil {
		return nil, err
	}

	dsn := sqliteDSN(absolutePath)

	database, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	database.SetMaxOpenConns(4)
	database.SetMaxIdleConns(4)
	database.SetConnMaxLifetime(0)
	database.SetConnMaxIdleTime(5 * time.Minute)

	if err := database.PingContext(ctx); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("ping sqlite database path=%s dsn=%s: %w", absolutePath, dsn, err)
	}

	if err := Migrate(ctx, database); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("migrate sqlite database: %w", err)
	}

	return database, nil
}

func sqliteDSN(databasePath string) string {
	p := filepath.ToSlash(databasePath)

	query := url.Values{}
	query.Set("mode", "rwc")

	query.Add("_pragma", "foreign_keys(1)")
	query.Add("_pragma", "busy_timeout(5000)")
	query.Add("_pragma", "journal_mode(WAL)")
	query.Add("_pragma", "synchronous(NORMAL)")

	return "file:" + p + "?" + query.Encode()
}

func ensureWritable(databasePath string) error {
	dir := filepath.Dir(databasePath)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create database directory %s: %w", dir, err)
	}

	testFile := filepath.Join(dir, ".write-test")

	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("database directory is not writable %s: %w", dir, err)
	}

	_ = f.Close()
	_ = os.Remove(testFile)

	return nil
}
