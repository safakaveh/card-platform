package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

var startTime = time.Now()

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// LiveHandler فقط بررسی می‌کند که خود اپلیکیشن بالا باشد
func (h *Handler) LiveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := Response{
			Status:    "ok",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Uptime:    int64(time.Since(h.service.startTime).Seconds()),
		}

		writeJSON(w, http.StatusOK, resp)
	}
}

// ReadyHandler وضعیت dependencyها را هم بررسی می‌کند
func (h *Handler) ReadyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checks := make(map[string]CheckResult)
		overallStatus := "ok"
		statusCode := http.StatusOK

		for _, checker := range h.service.checkers {
			ctx, cancel := context.WithTimeout(r.Context(), checker.Timeout)
			err := checker.Check(ctx)
			cancel()

			if err != nil {
				checks[checker.Name] = CheckResult{
					Status: "down",
					Error:  err.Error(),
				}
				overallStatus = "error"
				statusCode = http.StatusServiceUnavailable
			} else {
				checks[checker.Name] = CheckResult{
					Status: "ok",
				}
			}
		}

		resp := Response{
			Status:    overallStatus,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Uptime:    int64(time.Since(h.service.startTime).Seconds()),
			Checks:    checks,
		}

		writeJSON(w, statusCode, resp)
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(v)
}
