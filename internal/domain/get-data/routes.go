package getdata

import "github.com/go-chi/chi/v5"

func Routes(handler *Handler) chi.Router {
	r := chi.NewRouter()
	r.Get("/pending", handler.PendingMifare)
	r.Get("/read", handler.Read)
	r.Get("/read-report", handler.ReadReport)
	r.Post("/read/{cardID}/reset", handler.ResetRead)
	return r
}
