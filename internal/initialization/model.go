package initialization

import (
	"github.com/safakaveh/card-platform/internal/domain/health"
	"github.com/safakaveh/card-platform/internal/domain/shutdown"
)

type Model struct{}

type Handlers struct {
	HealthHandler   *health.Handler
	ShutdownHandler *shutdown.Handler
	// Application    *application.Handler
	// Auth           *auth.Handler
	// File           *file.Handler
	// RefreshToken   *refreshtoken.Handler
	// Rule           *rule.Handler
	// User           *user.Handler
	// UserPermission *userpermission.Handler
}

func NewHandler(ch *chan struct{}) *Handlers {
	return &Handlers{
		HealthHandler:   health.NewHandler(health.NewService()),
		ShutdownHandler: shutdown.NewHandler(shutdown.NewService(ch)),
	}
}
