package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	server *http.Server
}

func New(opts ...func(*Server)) *Server {
	srv := Server{
		server: &http.Server{
			Addr:         fmt.Sprintf(":%s", "8080"),
			IdleTimeout:  time.Minute,
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 10,
			Handler:      chi.NewMux(),
		},
	}

	for _, fn := range opts {
		fn(&srv)
	}

	return &srv
}

func WithRouter(router *chi.Mux) func(*Server) {
	return func(s *Server) {
		s.server.Handler = router
	}
}

func WithPort(port string) func(*Server) {
	return func(s *Server) {
		s.server.Addr = fmt.Sprintf(":%s", port)
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) GracefulShutdown(ctx context.Context, timeout time.Duration) chan error {
	shutdownErr := make(chan error)

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(
			stop,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)
		<-stop

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		shutdownErr <- s.Shutdown(ctx)
	}()

	return shutdownErr
}
