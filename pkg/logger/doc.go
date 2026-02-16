// Package logger provides a structured logger factory built on charmbracelet/log,
// with context propagation for automatic field enrichment across request lifecycles.
//
// The logger automatically switches format based on the environment:
// JSON output for production, human-readable text for development.
//
// # Creating a logger
//
// Basic usage with defaults (production, JSON, info level, stdout):
//
//	log := logger.New()
//
// Development logger with debug level:
//
//	log := logger.New(
//	    logger.WithEnvironment(logger.Development),
//	    logger.WithLevel(charmbraceletlog.DebugLevel),
//	    logger.WithPrefix("my-app"),
//	)
//
// # Context propagation
//
// Store an enriched logger in the request context via middleware, then retrieve
// it anywhere downstream. This avoids passing request-scoped fields (request ID,
// user ID, trace ID) to every log call manually.
//
// In your middleware, enrich the logger and store it in the context:
//
//	func RequestLogger(log *log.Logger) func(http.Handler) http.Handler {
//	    return func(next http.Handler) http.Handler {
//	        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	            enriched := log.With(
//	                "request_id", r.Header.Get("X-Request-ID"),
//	                "method", r.Method,
//	                "path", r.URL.Path,
//	            )
//	            ctx := logger.WithContext(r.Context(), enriched)
//	            next.ServeHTTP(w, r.WithContext(ctx))
//	        })
//	    }
//	}
//
// In your handlers and services, retrieve the logger from the context:
//
//	func (s service) CreateUser(ctx context.Context, name string) error {
//	    log := logger.FromContext(ctx)
//	    log.Info("creating user", "name", name)
//	    // request_id, method, path are included automatically
//	    return nil
//	}
//
// FromContext returns charmbracelet/log's default logger if no logger is found
// in the context, so it is always safe to call without nil checks.
//
// # Logging
//
//	log.Info("user created", "id", userID)
//	log.Error("payment failed", "err", err, "amount", amount)
//	log.Debug("cache hit", "key", cacheKey)
package logger
