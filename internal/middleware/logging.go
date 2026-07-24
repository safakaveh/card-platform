package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/safakaveh/card-platform/internal/common/logger"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logging اطلاعات مربوط به پردازش درخواست و پاسخ (مثل زمان اجرا و کد وضعیت) را ثبت می‌کند
func Logging(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{
				ResponseWriter: w,
				status:         200, // مقدار پیش‌فرض در صورت صدا زده نشدن WriteHeader
			}

			reqID := GetRequestID(r.Context())

			// افزودن جزئیات درخواست جاری به ساختار لاگر
			reqLogger := log.With(
				"request_id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
			)

			// ذخیره لاگر مجهز به فیلدهای جاری در Context
			ctx := logger.WithLogger(r.Context(), reqLogger)

			next.ServeHTTP(rw, r.WithContext(ctx))

			// ثبت خروجی پس از اتمام پردازش درخواست
			reqLogger.Info("request completed",
				"status", rw.status,
				"latency_ms", time.Since(start).Milliseconds(),
				"ip", r.RemoteAddr,
			)
		})
	}
}
