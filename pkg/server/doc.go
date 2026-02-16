// Package server provides an HTTP server wrapper with sensible defaults,
// chi router integration, and graceful shutdown via OS signals.
//
// The server uses the functional options pattern for configuration and comes
// with pre-configured timeouts (read: 5s, write: 10s, idle: 1m).
//
// Basic usage:
//
//	srv := server.New(
//	    server.WithRouter(router),
//	    server.WithPort("3000"),
//	)
//
//	shutdownErr := srv.GracefulShutdown(ctx, 30*time.Second)
//
//	if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
//	    log.Fatal("server failed", "err", err)
//	}
//
//	if err := <-shutdownErr; err != nil {
//	    log.Fatal("graceful shutdown failed", "err", err)
//	}
//
// GracefulShutdown listens for SIGHUP, SIGINT, SIGTERM, and SIGQUIT signals,
// then shuts down the server within the specified timeout, allowing in-flight
// requests to complete.
package server
