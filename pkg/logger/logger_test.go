package logger_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	charmlog "github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"

	"github.com/bernardinorafael/gogem/pkg/logger"
)

func TestNew_ReturnsLogger(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logger.New(logger.WithOutput(&buf))

	assert.NotNil(t, l)

	l.Info("smoke test")
	assert.NotEmpty(t, buf.String())
}

func TestNew_ProductionEnvironment_JSONFormat(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logger.New(
		logger.WithOutput(&buf),
		logger.WithEnvironment(logger.Production),
	)

	l.Info("test message")

	output := buf.String()
	assert.True(t, strings.HasPrefix(strings.TrimSpace(output), "{"), "production output should be JSON, got: %s", output)
	assert.Contains(t, output, "test message")
}

func TestNew_DevelopmentEnvironment_TextFormat(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logger.New(
		logger.WithOutput(&buf),
		logger.WithEnvironment(logger.Development),
	)

	l.Info("test message")

	output := buf.String()
	assert.False(t, strings.HasPrefix(strings.TrimSpace(output), "{"), "development output should not be JSON, got: %s", output)
	assert.Contains(t, output, "test message")
}

func TestNew_LevelFiltering(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		level       charmlog.Level
		logFn       func(l logger.Logger)
		expectEmpty bool
	}{
		{
			name:        "debug suppressed at info level",
			level:       charmlog.InfoLevel,
			logFn:       func(l logger.Logger) { l.Debug("hidden") },
			expectEmpty: true,
		},
		{
			name:        "debug passes at debug level",
			level:       charmlog.DebugLevel,
			logFn:       func(l logger.Logger) { l.Debug("visible") },
			expectEmpty: false,
		},
		{
			name:        "info passes at info level",
			level:       charmlog.InfoLevel,
			logFn:       func(l logger.Logger) { l.Info("visible") },
			expectEmpty: false,
		},
		{
			name:        "warn passes at info level",
			level:       charmlog.InfoLevel,
			logFn:       func(l logger.Logger) { l.Warn("visible") },
			expectEmpty: false,
		},
		{
			name:        "error passes at info level",
			level:       charmlog.InfoLevel,
			logFn:       func(l logger.Logger) { l.Error("visible") },
			expectEmpty: false,
		},
		{
			name:        "info suppressed at warn level",
			level:       charmlog.WarnLevel,
			logFn:       func(l logger.Logger) { l.Info("hidden") },
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			l := logger.New(
				logger.WithOutput(&buf),
				logger.WithLevel(tt.level),
			)

			tt.logFn(l)

			if tt.expectEmpty {
				assert.Empty(t, buf.String())
			} else {
				assert.NotEmpty(t, buf.String())
			}
		})
	}
}

func TestNew_WithPrefix_AppearsInOutput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logger.New(
		logger.WithOutput(&buf),
		logger.WithEnvironment(logger.Production),
		logger.WithPrefix("myapp"),
	)

	l.Info("hello")

	assert.Contains(t, buf.String(), "myapp")
}

func TestLogMethods_WriteCorrectOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		logFn  func(l logger.Logger)
		expect string
	}{
		{
			name:   "Info writes message",
			logFn:  func(l logger.Logger) { l.Info("info message") },
			expect: "info message",
		},
		{
			name:   "Infof writes formatted message",
			logFn:  func(l logger.Logger) { l.Infof("info %s", "formatted") },
			expect: "info formatted",
		},
		{
			name:   "Warn writes message",
			logFn:  func(l logger.Logger) { l.Warn("warn message") },
			expect: "warn message",
		},
		{
			name:   "Warnf writes formatted message",
			logFn:  func(l logger.Logger) { l.Warnf("warn %s", "formatted") },
			expect: "warn formatted",
		},
		{
			name:   "Error writes message",
			logFn:  func(l logger.Logger) { l.Error("error message") },
			expect: "error message",
		},
		{
			name:   "Errorf writes formatted message",
			logFn:  func(l logger.Logger) { l.Errorf("error %s", "formatted") },
			expect: "error formatted",
		},
		{
			name:   "Debug writes message at debug level",
			logFn:  func(l logger.Logger) { l.Debug("debug message") },
			expect: "debug message",
		},
		{
			name:   "Debugf writes formatted message at debug level",
			logFn:  func(l logger.Logger) { l.Debugf("debug %s", "formatted") },
			expect: "debug formatted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			l := logger.New(
				logger.WithOutput(&buf),
				logger.WithEnvironment(logger.Production),
				logger.WithLevel(charmlog.DebugLevel),
			)

			tt.logFn(l)

			assert.Contains(t, buf.String(), tt.expect)
		})
	}
}

func TestLogMethods_KeyValuePairs(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logger.New(
		logger.WithOutput(&buf),
		logger.WithEnvironment(logger.Production),
	)

	l.Info("user created", "id", "usr_123", "email", "test@example.com")

	output := buf.String()
	assert.Contains(t, output, "usr_123")
	assert.Contains(t, output, "test@example.com")
}

func TestWithContext_StoresAndRetrievesLogger(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logger.New(logger.WithOutput(&buf))

	ctx := logger.WithContext(context.Background(), l)
	retrieved := logger.FromContext(ctx)

	assert.NotNil(t, retrieved)

	retrieved.Info("from context")
	assert.Contains(t, buf.String(), "from context")
}

func TestFromContext_FallbackOnEmptyContext(t *testing.T) {
	t.Parallel()

	l := logger.FromContext(context.Background())
	assert.NotNil(t, l)
}

func TestWith_EnrichesOutput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logger.New(
		logger.WithOutput(&buf),
		logger.WithEnvironment(logger.Production),
	)

	enriched := l.With("request_id", "req_abc123")
	enriched.Info("user action")

	output := buf.String()
	assert.Contains(t, output, "req_abc123")
	assert.Contains(t, output, "user action")
}

func TestWith_DoesNotAffectOriginalLogger(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logger.New(
		logger.WithOutput(&buf),
		logger.WithEnvironment(logger.Production),
	)

	_ = l.With("request_id", "req_abc123")

	l.Info("original")

	assert.NotContains(t, buf.String(), "req_abc123")
	assert.Contains(t, buf.String(), "original")
}

func TestWith_ReturnsLogger(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logger.New(
		logger.WithOutput(&buf),
		logger.WithEnvironment(logger.Production),
	)

	enriched := l.With("key", "value")
	assert.NotNil(t, enriched)

	var _ logger.Logger = enriched
}

func TestFromContext_NestedContext(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	l := logger.New(logger.WithOutput(&buf))

	parent := logger.WithContext(context.Background(), l)
	child, cancel := context.WithCancel(parent)
	defer cancel()

	retrieved := logger.FromContext(child)
	assert.NotNil(t, retrieved)

	retrieved.Info("from child context")
	assert.Contains(t, buf.String(), "from child context")
}
