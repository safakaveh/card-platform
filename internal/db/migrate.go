package db

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

const createMigrationsTableSQL = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT PRIMARY KEY NOT NULL,
    checksum TEXT NOT NULL,
    applied_at INTEGER NOT NULL
);
`

func Migrate(ctx context.Context, database *sql.DB) error {
	if database == nil {
		return errors.New("database is nil")
	}

	if _, err := database.ExecContext(ctx, createMigrationsTableSQL); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	entries, err := fs.ReadDir(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("read embedded migrations: %w", err)
	}

	migrations := make([]string, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()

		if strings.HasSuffix(strings.ToLower(name), ".up.sql") {
			migrations = append(migrations, name)
		}
	}

	sort.Strings(migrations)

	for _, name := range migrations {
		if err := applyMigration(ctx, database, name); err != nil {
			return err
		}
	}

	return nil
}

func applyMigration(
	ctx context.Context,
	database *sql.DB,
	name string,
) error {
	path := "migrations/" + name

	content, err := migrationFiles.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read migration %q: %w", name, err)
	}

	checksum := calculateChecksum(content)

	appliedChecksum, applied, err := findAppliedMigration(
		ctx,
		database,
		name,
	)
	if err != nil {
		return fmt.Errorf("check migration %q: %w", name, err)
	}

	if applied {
		if appliedChecksum != checksum {
			return fmt.Errorf(
				"migration %q was modified after being applied: stored checksum=%s, current checksum=%s",
				name,
				appliedChecksum,
				checksum,
			)
		}

		return nil
	}

	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration %q: %w", name, err)
	}

	committed := false

	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	if _, err := tx.ExecContext(ctx, string(content)); err != nil {
		return fmt.Errorf("execute migration %q: %w", name, err)
	}

	const insertMigrationSQL = `
		INSERT INTO schema_migrations (
			version,
			checksum,
			applied_at
		)
		VALUES (
			?,
			?,
			CAST((julianday('now') - 2440587.5) * 86400000 AS INTEGER)
		);
	`

	if _, err := tx.ExecContext(
		ctx,
		insertMigrationSQL,
		name,
		checksum,
	); err != nil {
		return fmt.Errorf("record migration %q: %w", name, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %q: %w", name, err)
	}

	committed = true

	return nil
}

func findAppliedMigration(
	ctx context.Context,
	database *sql.DB,
	version string,
) (
	checksum string,
	applied bool,
	err error,
) {
	const query = `
		SELECT checksum
		FROM schema_migrations
		WHERE version = ?;
	`

	err = database.QueryRowContext(ctx, query, version).Scan(&checksum)

	switch {
	case err == nil:
		return checksum, true, nil

	case errors.Is(err, sql.ErrNoRows):
		return "", false, nil

	default:
		return "", false, err
	}
}

func calculateChecksum(content []byte) string {
	sum := sha256.Sum256(content)

	return hex.EncodeToString(sum[:])
}
