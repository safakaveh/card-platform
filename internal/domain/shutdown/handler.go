package shutdown

import (
	"net"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) ShutDown(w http.ResponseWriter, r *http.Request) {
	if !isLoopbackRequest(r) {
		http.Error(
			w,
			"shutdown is only available locally",
			http.StatusForbidden,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(
		`{"message":"application is shutting down"}`,
	))

	h.service.Shutdown()
}

func isLoopbackRequest(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}

	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
