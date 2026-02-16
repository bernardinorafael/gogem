package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

type loggerKey struct{}

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

func New(opts ...func(*config)) *log.Logger {
	cfg := config{
		output:      os.Stdout,
		level:       log.InfoLevel,
		environment: Production,
		timeFormat:  time.RFC3339,
	}

	for _, fn := range opts {
		fn(&cfg)
	}

	formatter := log.JSONFormatter
	timeFormat := cfg.timeFormat

	if cfg.environment == Development {
		formatter = log.TextFormatter
		timeFormat = time.Kitchen
	}

	logger := log.NewWithOptions(cfg.output, log.Options{
		TimeFormat:      timeFormat,
		Formatter:       formatter,
		ReportTimestamp: true,
		Prefix:          cfg.prefix,
		Level:           cfg.level,
	})

	return logger
}

type config struct {
	output      io.Writer
	level       log.Level
	environment Environment
	prefix      string
	timeFormat  string
}

func WithOutput(w io.Writer) func(*config) {
	return func(c *config) {
		c.output = w
	}
}

func WithLevel(level log.Level) func(*config) {
	return func(c *config) {
		c.level = level
	}
}

func WithEnvironment(env Environment) func(*config) {
	return func(c *config) {
		c.environment = env
	}
}

func WithPrefix(prefix string) func(*config) {
	return func(c *config) {
		c.prefix = prefix
	}
}

func WithTimeFormat(format string) func(*config) {
	return func(c *config) {
		c.timeFormat = format
	}
}

// WithContext stores a logger in the context. Use this in middleware
// to propagate an enriched logger through the request lifecycle.
func WithContext(ctx context.Context, l *log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

// FromContext retrieves the logger from the context.
// Returns the default logger (production, info level) if none is found.
func FromContext(ctx context.Context) *log.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*log.Logger); ok {
		return l
	}
	return log.Default()
}
