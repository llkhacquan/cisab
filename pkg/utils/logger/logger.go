package logger

import (
	"io"
	"log/slog"
	"os"
)

// Logger is a wrapper around slog.Logger
type Logger struct {
	*slog.Logger
}

// Config holds configuration for the logger
type Config struct {
	// Level is the minimum log level that will be logged
	Level slog.Level
	// Output is where the logs will be written to
	Output io.Writer
	// AddSource adds the file name and line number to log messages
	AddSource bool
}

// DefaultConfig returns a default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:     slog.LevelInfo,
		Output:    os.Stdout,
		AddSource: false,
	}
}

// New creates a new Logger with the given configuration
func New(cfg Config) *Logger {
	var handler slog.Handler

	// Create a default JSON handler
	handlerOptions := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}
	handler = slog.NewJSONHandler(cfg.Output, handlerOptions)

	return &Logger{
		Logger: slog.New(handler),
	}
}

// NewDefault creates a new Logger with default configuration
func NewDefault() *Logger {
	return New(DefaultConfig())
}

// With returns a new Logger with the given attributes added to it
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
	}
}
