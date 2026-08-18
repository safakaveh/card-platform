package getdata

import "github.com/go-chi/chi/v5"

func Routes(handler *Handler) chi.Router {
	r := chi.NewRouter()
	r.Get("/pending", handler.PendingMifare)
	return r
}
