// Package logger provides a structured logger factory built on charmbracelet/log.
//
// The logger automatically switches format based on the environment:
// JSON output for production, human-readable text for development.
//
// Creating a logger:
//
//	log := logger.NewLogger(true, "my-app", "development")
//	// Output: 3:04PM INF my-app: user created id=user_abc123
//
//	log := logger.NewLogger(false, "my-app", "production")
//	// Output: {"time":"2024-01-15T10:30:00Z","level":"info","prefix":"my-app","msg":"user created","id":"user_abc123"}
//
// Usage:
//
//	log.Info("user created", "id", userID)
//	log.Error("payment failed", "err", err, "amount", amount)
//	log.Debug("cache hit", "key", cacheKey) // only visible when debug=true
package logger
