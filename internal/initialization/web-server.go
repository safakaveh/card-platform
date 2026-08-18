package initialization

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/safakaveh/card-platform/internal/config"
)

type Server struct {
	Addr            string
	Handler         http.Handler
	shutdownRequest chan struct{}
}

func NewWebServer(r http.Handler, ch chan struct{}) *Server {
	port := config.GetEnvConf().AppHttpPort
	return &Server{
		Addr:            fmt.Sprintf("127.0.0.1:%d", port),
		Handler:         r,
		shutdownRequest: ch,
	}
}

func (s Server) Start() {
	server := &http.Server{
		Addr:              s.Addr,
		Handler:           s.Handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       0,
		WriteTimeout:      0,
		IdleTimeout:       60 * time.Second,
	}
	url := "http://" + s.Addr

	// Bind synchronously so a second launch is detected deterministically.
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		if notifyDuplicateInstance(url) {
			_ = s.openBrowser(url)
		}
		return
	}

	serverErr := make(chan error, 1)
	signalCtx, stopSignals := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stopSignals()

	go func() {
		if err := s.openBrowser(url); err != nil {
			log.Printf("failed to open browser automatically: %v", err)
			log.Printf("open manually: %s", url)
		}
	}()
	go func() {

		log.Printf("server listening on %s", url)

		err := server.Serve(listener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}

		serverErr <- nil
	}()

	select {
	case <-s.shutdownRequest:
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

func (s Server) openBrowser(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("rundll32", "url.dll", "FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}
