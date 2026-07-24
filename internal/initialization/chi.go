package initialization

import (
	"errors"
	"io/fs"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/safakaveh/card-platform/internal/common/logger"
	"github.com/safakaveh/card-platform/internal/domain/health"
	"github.com/safakaveh/card-platform/internal/domain/shutdown"
	"github.com/safakaveh/card-platform/internal/middleware"
	"github.com/safakaveh/card-platform/internal/web"
)

type ChiRouter struct {
	Router *chi.Mux
}

func NewChiRoute(h *Handlers) *ChiRouter {
	log := logger.New()
	r := chi.NewRouter()

	// middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.Recovery)
	r.Use(middleware.Logging(log))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "QUERY"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/health", func(api chi.Router) {
		api.Mount("/", health.Routes(h.HealthHandler))
	})
	r.Route("/system", func(api chi.Router) {
		api.Mount("/shutdown", shutdown.Routes(h.ShutdownHandler))
	})

	// r.Group(func(api chi.Router) {
	// 	api.Use(jwtauth.AuthMiddleware(v))

	// 	api.Route("/applications", func(sub chi.Router) {
	// 		sub.Mount("/", application.Routes(h.Application))
	// 	})
	// 	api.Route("/files", func(sub chi.Router) {
	// 		sub.Mount("/", file.Routes(h.File))
	// 	})
	// 	api.Route("/refreshtokens", func(sub chi.Router) {
	// 		sub.Mount("/", refreshtoken.Routes(h.RefreshToken))
	// 	})
	// 	api.Route("/rules", func(sub chi.Router) {
	// 		sub.Mount("/", rule.Routes(h.Rule))
	// 	})
	// 	api.Route("/users", func(sub chi.Router) {
	// 		sub.Mount("/", user.Routes(h.User))
	// 	})
	// 	api.Route("/userpermissions", func(sub chi.Router) {
	// 		sub.Mount("/", userpermission.Routes(h.UserPermission))
	// 	})
	// })

	distFS, err := fs.Sub(web.Assets, "build")
	if err != nil {
		log.Error(err.Error())
	}

	r.NotFound(spaHandler(distFS))

	return &ChiRouter{
		Router: r,
	}
}

func spaHandler(distFS fs.FS) http.HandlerFunc {
	fileServer := http.FileServer(http.FS(distFS))

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.NotFound(w, r)
			return
		}

		p := r.URL.Path
		if p == "/" {
			p = "/index.html"
		}

		_, err := fs.Stat(distFS, p[1:])
		if err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}
		if !errors.Is(err, os.ErrNotExist) {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		r2 := r.Clone(r.Context())
		r2.URL.Path = "/index.html"
		fileServer.ServeHTTP(w, r2)
	}
}
