// Package server provides an HTTP server wrapper with sensible defaults
// and graceful shutdown via OS signals.
//
// The server uses the functional options pattern for configuration and comes
// with pre-configured timeouts (read: 5s, write: 10s, idle: 1m, shutdown: 30s).
//
// Basic usage:
//
//	srv := server.New(
//	    server.WithHandler(router),
//	    server.WithPort(3000),
//	)
//
//	if err := srv.ListenAndServe(); err != nil {
//	    log.Fatal("server error", "err", err)
//	}
//
// ListenAndServe blocks until a shutdown signal is received (SIGHUP, SIGINT,
// SIGTERM, SIGQUIT), then gracefully shuts down allowing in-flight requests
// to complete within the configured timeout.
//
// Custom configuration:
//
//	srv := server.New(
//	    server.WithHandler(router),
//	    server.WithAddr("127.0.0.1:3000"),
//	    server.WithReadTimeout(10*time.Second),
//	    server.WithWriteTimeout(30*time.Second),
//	    server.WithShutdownTimeout(60*time.Second),
//	)
//
// WithHandler accepts any http.Handler, so it works with chi, gorilla, stdlib mux,
// or any custom router.
package server
