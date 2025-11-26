package humanlog

import (
	"io"
	"log/slog"
)

// Options configures the human-readable slog.Handler.
type Options struct {
	// Level is the minimum level to log.
	Level slog.Level

	// Writer is where the logs are written to.
	Writer io.Writer

	// TimeFormat is the format used for timestamps.
	// Default: "15:04:05" (hour:minute:second)
	TimeFormat string

	// DisableColor disables colored output for log levels.
	// When true, no ANSI color codes will be used.
	DisableColor bool

	// AddSource causes the handler to compute the source code position
	// of the log statement and add a "source" attribute to the output.
	AddSource bool

	// MessageWidth sets the fixed width for the message field.
	// Messages longer than this width will be truncated with "...".
	// Messages shorter will be padded with spaces.
	// Default: 40 characters
	MessageWidth int

	// UseJSON enables JSON output format instead of human-readable format.
	// When true, the handler will delegate to slog.JSONHandler for structured output.
	// Useful for production environments where log aggregation systems expect JSON.
	UseJSON bool
}

// DefaultOptions returns a new Options with default values.
func DefaultOptions() *Options {
	return &Options{
		Level:        slog.LevelInfo,
		TimeFormat:   "15:04:05",
		DisableColor: false,
		AddSource:    true,
		MessageWidth: 40,
		UseJSON:      false,
	}
}
