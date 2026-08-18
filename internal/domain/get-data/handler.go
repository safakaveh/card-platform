package getdata

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }

func (h *Handler) PendingMifare(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	result, err := h.service.PendingMifare(r.Context(), limit)
	if err != nil {
		http.Error(w, `{"error":"دریافت داده‌های خوانده‌نشده ناموفق بود"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = writeJSON(w, result)
}

func writeJSON(w http.ResponseWriter, value any) error {
	return json.NewEncoder(w).Encode(value)
}
