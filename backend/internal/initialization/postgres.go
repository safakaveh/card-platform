package initialization

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sunecity/smart-building-platform/auth/internal/config"
)

type PostgresData struct {
	Pool *pgxpool.Pool
}

func newPostgres() (*PostgresData, error) {
	conf := config.GetEnvConf()

	pool, err := pgxpool.New(context.Background(), conf.DbDsn)
	if err != nil {
		return nil, ErrPostgresInit.WithCause(err).LogOnly()
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, ErrPostgresInit.WithCause(err).LogOnly()
	}

	return &PostgresData{Pool: pool}, nil
}

var (
	postgresData *PostgresData
	postgresMu   sync.Mutex
)

func GetPostgresData() (*PostgresData, error) {
	postgresMu.Lock()
	defer postgresMu.Unlock()

	if postgresData != nil {
		return postgresData, nil
	}

	pg, err := newPostgres()
	if err != nil {
		return nil, err
	}

	postgresData = pg
	return postgresData, nil
}

func ClosePostgres() {
	postgresMu.Lock()
	defer postgresMu.Unlock()

	if postgresData != nil && postgresData.Pool != nil {
		postgresData.Pool.Close()
		postgresData = nil
	}
}
