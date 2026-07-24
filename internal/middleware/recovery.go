package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/safakaveh/card-platform/internal/common/logger"
)

// Recovery جلوی کرش کردن برنامه در اثر خطاهای ناگهانی (Panic) را می‌گیرد و آن را لاگ می‌کند
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				l := logger.FromContext(r.Context())

				l.Error("panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
				)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(`{"error":"internal server error"}`))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
