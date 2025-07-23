package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/lepinkainen/humanlog"
)

func main() {
	// --- Default Handler ---
	// Create a new human-readable handler with default options.
	// By default, AddSource is true, which adds the source file and line number.
	defaultHandler := humanlog.NewHandler(os.Stdout, nil) // nil opts -> use defaults
	logger := slog.New(defaultHandler)

	// Set the default logger to see output from top-level slog functions.
	slog.SetDefault(logger)

	slog.Info("This message comes from the default logger.")
	slog.Info("This is a message that will be truncated because it's too long for the fixed-width field, which is a feature of this handler.")

	// --- Logging with different levels and attributes ---
	logger.Debug("Debug message with attributes", slog.Int("count", 42), slog.Bool("enabled", true))
	logger.Info("Info message with attributes", slog.String("user", "john"), slog.String("action", "login"))
	logger.Warn("Warning message with attributes", slog.Duration("latency", 150*time.Millisecond), slog.Duration("threshold", 100*time.Millisecond))
	logger.Error("Error message with attributes", slog.String("error", "connection refused"), slog.Int("retries", 3))

	// --- Structured logging with LogAttrs ---
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Structured attributes",
		slog.String("service", "auth"),
		slog.Int("status", 200),
		slog.Duration("response_time", 30*time.Millisecond),
	)

	// --- Logging with Groups ---
	// Attributes logged with this logger will be prefixed with the group name.
	groupedLogger := logger.WithGroup("request")
	groupedLogger.Info("Request received",
		slog.String("method", "GET"),
		slog.String("path", "/api/v1/users"),
		slog.String("remote_addr", "192.168.1.1"),
	)

	// --- Logging with pre-registered attributes ---
	// These attributes are included in all subsequent logs from this logger instance.
	requestLogger := logger.With(slog.String("request_id", "abc-123"), slog.String("trace_id", "xyz-789"))
	requestLogger.Info("Processing request")
	requestLogger.Error("Request failed", slog.String("error", "database connection error"))

	// --- Combining Groups and pre-registered attributes ---
	// Attributes are prefixed with the group name.
	userRequestLogger := logger.WithGroup("user_request").With(slog.Int("user_id", 42))
	userRequestLogger.Info("User action", slog.String("action", "update_profile"))

	// --- Custom Handler without source ---
	slog.Info("--- Now logging with a custom handler (AddSource disabled) ---")
	customOpts := &humanlog.Options{
		Level:      slog.LevelInfo,
		TimeFormat: time.Kitchen,
		AddSource:  false, // Explicitly disable source location
	}
	handlerNoSource := humanlog.NewHandler(os.Stdout, customOpts)
	loggerNoSource := slog.New(handlerNoSource)
	loggerNoSource.Info("This log comes from a custom handler and does not have source information.")
	loggerNoSource.Warn("Note the different time format.")
}
