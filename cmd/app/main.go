package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/safakaveh/card-platform/internal/db"
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
	port := "8080"
	addr := "127.0.0.1:" + port
	url := "http://" + addr

	sqliteDB, err := db.Open(context.Background(), "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := sqliteDB.Close(); err != nil {
			log.Printf("failed to close database: %v", err)
		} else {
			log.Println("database connection closed")
		}
	}()

	srv := &Server{
		rawDB:   sqliteDB,
		queries: sqlc.New(sqliteDB),
	}

	shutdownRequest := make(chan struct{}, 1)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Post("/api/system/shutdown", func(w http.ResponseWriter, r *http.Request) {
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

		select {
		case shutdownRequest <- struct{}{}:
		default:
		}
	})

	r.Route("/api/users", func(r chi.Router) {

		r.Post("/", srv.handleCreateUser)
		r.Delete("/{id}", srv.handleDeleteUser)
	})

	distFS, err := fs.Sub(web.Assets, "dist")
	if err != nil {
		log.Fatal(err)
	}

	r.NotFound(spaHandler(distFS))

	server := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		if err := waitForServer(addr, 10*time.Second); err != nil {
			log.Printf("server not ready: %v", err)
			log.Printf("open manually: %s", url)
			return
		}

		if err := openBrowser(url); err != nil {
			log.Printf("failed to open browser automatically: %v", err)
			log.Printf("open manually: %s", url)
		}
	}()

	signalCtx, stopSignals := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stopSignals()

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("server listening on %s", url)

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}

		serverErr <- nil
	}()

	select {
	case <-shutdownRequest:
		log.Println("shutdown requested through Web API")

	case <-signalCtx.Done():
		log.Println("shutdown requested through OS signal")

	case err := <-serverErr:
		if err != nil {
			log.Printf("HTTP server failed: %v", err)
		}
		return
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancelShutdown()

	log.Println("shutting down HTTP server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)

		if closeErr := server.Close(); closeErr != nil {
			log.Printf("force close failed: %v", closeErr)
		}
	}

	log.Println("application stopped")
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

func waitForServer(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 300*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		time.Sleep(150 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for %s", addr)
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll", "FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}

func isLoopbackRequest(r *http.Request) bool {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false
	}

	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
