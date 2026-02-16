package logger

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func NewLogger(debug bool, appName, environment string) *log.Logger {
	logger := log.NewWithOptions(
		os.Stdout,
		log.Options{
			TimeFormat:      time.RFC3339,
			Formatter:       log.JSONFormatter,
			ReportTimestamp: true,
			Prefix:          appName,
		},
	)

	if debug {
		logger.SetLevel(log.DebugLevel)
	}

	// In development environment, use TextFormatter for easier reading
	if environment == "development" {
		logger.SetFormatter(log.TextFormatter)
		logger.SetTimeFormat(time.Kitchen)
	}

	return logger
}
