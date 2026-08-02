package uploadcsv

import "github.com/go-chi/chi/v5"

func Routes(handler *Handler) chi.Router {
	router := chi.NewRouter()
	router.Post("/", handler.Upload)
	router.Get("/", handler.List)
	router.Get("/{id}", handler.Get)
	router.Delete("/{id}", handler.Delete)
	return router
}
