package initialization

import (
	"context"
	"database/sql"
	"log"

	"github.com/safakaveh/card-platform/internal/config"
	"github.com/safakaveh/card-platform/internal/db"
	"github.com/safakaveh/card-platform/internal/db/sqlc"
)

type DBServer struct {
	Database *sql.DB
	Queries  *sqlc.Queries
}

func NewDBServer() *DBServer {
	sqliteDB, err := db.Open(context.Background(), config.GetEnvConf().DbDatasourceName)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := sqliteDB.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		} else {
			log.Println("database connection closed")
		}
	}()
	return &DBServer{
		Database: sqliteDB,
		Queries:  sqlc.New(sqliteDB),
	}
}
