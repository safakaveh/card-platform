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
		log.Printf("cannot start application on %s: %v", s.Addr, err)
		if notifyDuplicateInstance(url) {
			if browserErr := s.openBrowser(url); browserErr != nil {
				log.Printf("failed to open the running application: %v", browserErr)
				log.Printf("open manually: %s", url)
			}
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
	command, args, err := browserCommand(runtime.GOOS, url)
	if err != nil {
		return err
	}
	return exec.Command(command, args...).Start()
}

// browserCommand returns a platform-native command without assuming that the
// executable is installed. The fallback commands make Linux distributions
// using gio/sensible-browser work even when xdg-open is unavailable.
func browserCommand(goos, url string) (string, []string, error) {
	switch goos {
	case "windows":
		return "rundll32", []string{"url.dll", "FileProtocolHandler", url}, nil
	case "darwin":
		return "open", []string{url}, nil
	default:
		for _, candidate := range []struct {
			name string
			args []string
		}{
			{"xdg-open", []string{url}},
			{"gio", []string{"open", url}},
			{"sensible-browser", []string{url}},
		} {
			if _, err := exec.LookPath(candidate.name); err == nil {
				return candidate.name, candidate.args, nil
			}
		}
		return "", nil, fmt.Errorf("no supported browser opener found (tried xdg-open, gio, sensible-browser)")
	}
}
