package health

import "github.com/go-chi/chi/v5"

func Routes(h *Handler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.LiveHandler())
	r.Get("/liveness", h.LiveHandler())
	r.Get("/readiness", h.ReadyHandler())

	return r
}
