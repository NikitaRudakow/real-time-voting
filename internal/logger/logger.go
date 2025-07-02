package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func New(level string) zerolog.Logger {
	// Parse log level
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	// Configure zerolog
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Create logger with console writer for development
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}
