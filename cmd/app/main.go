package main

import (
	"database/sql"
	"log"

	"github.com/safakaveh/card-platform/internal/db/sqlc"
	"github.com/safakaveh/card-platform/internal/initialization"
	_ "modernc.org/sqlite"
)

type Server struct {
	rawDB   *sql.DB
	queries *sqlc.Queries
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

func main() {
	shutdownRequest := make(chan struct{}, 1)

	DBServer := initialization.NewDBServer()
	defer func() {
		if err := DBServer.Database.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		} else {
			log.Println("database connection closed")
		}
	}()

	handlers := initialization.NewHandler(&shutdownRequest)
	chiRouters := initialization.NewChiRoute(handlers)

	server := initialization.NewWebServer(chiRouters.Router, shutdownRequest)
	server.Start()
}
