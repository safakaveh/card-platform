package initialization

import (
	"database/sql"

	getdata "github.com/safakaveh/card-platform/internal/domain/get-data"
	"github.com/safakaveh/card-platform/internal/domain/health"
	"github.com/safakaveh/card-platform/internal/domain/shutdown"
	uploadcsv "github.com/safakaveh/card-platform/internal/domain/upload-csv"
)

type Model struct{}

type Handlers struct {
	HealthHandler    *health.Handler
	ShutdownHandler  *shutdown.Handler
	UploadCSVHandler *uploadcsv.Handler
	GetDataHandler   *getdata.Handler
	// Application    *application.Handler
	// Auth           *auth.Handler
	// File           *file.Handler
	// RefreshToken   *refreshtoken.Handler
	// Rule           *rule.Handler
	// User           *user.Handler
	// UserPermission *userpermission.Handler
}

func NewHandler(ch *chan struct{}, database *sql.DB) *Handlers {
	return &Handlers{
		HealthHandler:    health.NewHandler(health.NewService()),
		ShutdownHandler:  shutdown.NewHandler(shutdown.NewService(ch)),
		UploadCSVHandler: uploadcsv.NewHandler(uploadcsv.NewService(database)),
		GetDataHandler:   getdata.NewHandler(getdata.NewService(database)),
	}
}
