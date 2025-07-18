// Package humanlog provides a human-readable formatter for Go's log/slog package.
//
// It implements slog.Handler to produce logs that are easier for humans to scan and interpret.
//
// Features:
//   - Fixed-width, colorized log level and message formatting
//   - Configurable log level, color output, and time format
//   - Drop-in replacement for slog.Handler
//
// Example usage:
//
//	handler := humanlog.NewHandler(os.Stdout, nil)
//	logger := slog.New(handler)
//	logger.Info("Hello, human-readable logs!")
//
// For advanced configuration and usage patterns, see example/main.go.
package humanlog
