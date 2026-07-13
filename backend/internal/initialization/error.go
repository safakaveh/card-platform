package initialization

import (
	"net/http"

	"github.com/sunecity/smart-building-platform/auth/internal/common/microerror"
)

var (
	ErrNatsInit     = microerror.New("init_01", http.StatusServiceUnavailable, "nats initialaization has error")
	ErrPostgresInit = microerror.New("init_02", http.StatusServiceUnavailable, "postgres initialaization has error")
	ErrVaultInit    = microerror.New("init_03", http.StatusServiceUnavailable, "vault initialaization has error")
)
