package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

func New(opts ...func(*Server)) *Server {
	srv := Server{
		server: &http.Server{
			Addr:         ":8080",
			IdleTimeout:  time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			Handler:      http.DefaultServeMux,
		},
		shutdownTimeout: 30 * time.Second,
	}

	for _, fn := range opts {
		fn(&srv)
	}

	return &srv
}

func WithHandler(handler http.Handler) func(*Server) {
	return func(s *Server) {
		s.server.Handler = handler
	}
}

func WithPort(port int) func(*Server) {
	return func(s *Server) {
		s.server.Addr = fmt.Sprintf(":%d", port)
	}
}

func WithAddr(addr string) func(*Server) {
	return func(s *Server) {
		s.server.Addr = addr
	}
}

func WithReadTimeout(d time.Duration) func(*Server) {
	return func(s *Server) {
		s.server.ReadTimeout = d
	}
}

func WithWriteTimeout(d time.Duration) func(*Server) {
	return func(s *Server) {
		s.server.WriteTimeout = d
	}
}

func WithIdleTimeout(d time.Duration) func(*Server) {
	return func(s *Server) {
		s.server.IdleTimeout = d
	}
}

func WithShutdownTimeout(d time.Duration) func(*Server) {
	return func(s *Server) {
		s.shutdownTimeout = d
	}
}

// ListenAndServe starts the HTTP server and blocks until a shutdown signal
// is received (SIGHUP, SIGINT, SIGTERM, SIGQUIT). It then gracefully shuts
// down the server, allowing in-flight requests to complete within the
// configured shutdown timeout.
func (s *Server) ListenAndServe() error {
	errCh := make(chan error, 1)

	go func() {
		if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case err := <-errCh:
		return err
	case <-stop:
		ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()
		return s.server.Shutdown(ctx)
	}
}

// Addr returns the server's listen address.
func (s *Server) Addr() string {
	return s.server.Addr
}
