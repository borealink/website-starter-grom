package logging

import (
	"io"
	"log"
)

// Logger wraps the standard Go logger.
type Logger struct {
	logger *log.Logger
}

// Config contains logger configuration options.
type Config struct {

	// Output destination for log messages.
	Out io.Writer

	// Prefix added before each log message.
	Prefix string

	// Formatting flags used by the logger.
	Flags int
}

// NewLogger creates and initializes a new logger instance.
func NewLogger(
	config Config,
) *Logger {
	return &Logger{
		logger: log.New(
			config.Out,
			config.Prefix,
			config.Flags,
		),
	}
}

// Get returns the underlying standard logger instance.
func (logger *Logger) Get() *log.Logger {
	return logger.logger
}
