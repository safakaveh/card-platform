package getdata

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	orderName := r.URL.Query().Get("order_name")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	result, err := h.service.ReadNext(r.Context(), orderName, limit)
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrOrderNameRequired || err == ErrOrderNotFound {
			status = http.StatusBadRequest
		}
		writeError(w, status, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = writeJSON(w, result)
}

func (h *Handler) ReadReport(w http.ResponseWriter, r *http.Request) {
	orderName := r.URL.Query().Get("order_name")
	status := r.URL.Query().Get("status")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	result, err := h.service.ReadReport(r.Context(), orderName, status, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = writeJSON(w, result)
}

func (h *Handler) ResetRead(w http.ResponseWriter, r *http.Request) {
	if err := h.service.ResetRead(r.Context(), chi.URLParam(r, "cardID")); err != nil {
		if err == ErrCardNotFound {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, value any) error {
	return json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
