package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/safakaveh/card-platform/internal/db/sqlc"
	"github.com/safakaveh/card-platform/internal/web"
	_ "modernc.org/sqlite"
)

type Server struct {
	rawDB   *sql.DB
	queries *sqlc.Queries
}

type CreateUserRequest struct {
	Name string `json:"name"`
}

func main() {
	sqliteDB, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer sqliteDB.Close()

	_, err = sqliteDB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	srv := &Server{
		rawDB:   sqliteDB,
		queries: sqlc.New(sqliteDB),
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/users", func(r chi.Router) {

		r.Post("/", srv.handleCreateUser)
		r.Delete("/{id}", srv.handleDeleteUser)
	})

	distFS, err := fs.Sub(web.Assets, "dist")
	if err != nil {
		log.Fatal(err)
	}

	r.NotFound(spaHandler(distFS))

	log.Println("open http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	res, err := s.rawDB.ExecContext(r.Context(),
		`INSERT INTO users (name) VALUES (?)`,
		req.Name,
	)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	id, _ := res.LastInsertId()

	writeJSON(w, map[string]any{
		"id":   id,
		"name": req.Name,
	})
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	_, err = s.rawDB.ExecContext(r.Context(),
		`DELETE FROM users WHERE id = ?`,
		id,
	)
	if err != nil {
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
