package initialization

import (
	"database/sql"

	"github.com/safakaveh/card-platform/internal/db/sqlc"
)

type Server struct {
	rawDB   *sql.DB
	queries *sqlc.Queries
}
