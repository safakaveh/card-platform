package initialization

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/safakaveh/card-platform/internal/config"
	"github.com/safakaveh/card-platform/internal/db"
	"github.com/safakaveh/card-platform/internal/db/sqlc"
)

type DBServer struct {
	Database *sql.DB
	Queries  *sqlc.Queries
}

func NewDBServer() *DBServer {
	datasource := config.GetEnvConf().DbDatasourceName
	if runtime.GOOS == "windows" && !filepath.IsAbs(datasource) {
		// A GUI-launched Windows executable does not always inherit the
		// executable directory as its working directory.
		if executable, err := os.Executable(); err == nil {
			datasource = filepath.Join(filepath.Dir(executable), datasource)
		}
	}

	sqliteDB, err := db.Open(context.Background(), datasource)
	if err != nil {
		log.Fatal(err)
	}

	return &DBServer{
		Database: sqliteDB,
		Queries:  sqlc.New(sqliteDB),
	}
}
