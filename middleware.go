package humanlog

import (
	"context"
	"log/slog"
)

// ContextKey is a type for context keys to avoid collisions
type ContextKey string

const (
	// RequestIDKey is the context key for request IDs
	RequestIDKey ContextKey = "request_id"
	// TraceIDKey is the context key for trace IDs
	TraceIDKey ContextKey = "trace_id"
	// UserIDKey is the context key for user IDs
	UserIDKey ContextKey = "user_id"
)

// WithRequestID adds a request ID to the context that will be automatically
// included in all log entries when using the contextual logging functions.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithTraceID adds a trace ID to the context that will be automatically
// included in all log entries when using the contextual logging functions.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithUserID adds a user ID to the context that will be automatically
// included in all log entries when using the contextual logging functions.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// ContextLogger wraps an slog.Logger to automatically extract and include
// correlation IDs from context in log entries.
type ContextLogger struct {
	logger *slog.Logger
}

// NewContextLogger creates a new context-aware logger that automatically
// includes correlation IDs from the context.
func NewContextLogger(logger *slog.Logger) *ContextLogger {
	return &ContextLogger{logger: logger}
}

// extractContextAttrs extracts correlation attributes from context
func (cl *ContextLogger) extractContextAttrs(ctx context.Context) []slog.Attr {
	var attrs []slog.Attr

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok && requestID != "" {
		attrs = append(attrs, slog.String("request_id", requestID))
	}

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		attrs = append(attrs, slog.String("trace_id", traceID))
	}

	if userID, ok := ctx.Value(UserIDKey).(string); ok && userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}

	return attrs
}

// LogAttrs logs with both context-extracted attributes and provided attributes
func (cl *ContextLogger) LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	contextAttrs := cl.extractContextAttrs(ctx)
	allAttrs := append(contextAttrs, attrs...)
	cl.logger.LogAttrs(ctx, level, msg, allAttrs...)
}

// Info logs an info message with context attributes
func (cl *ContextLogger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	cl.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
}

// Error logs an error message with context attributes
func (cl *ContextLogger) Error(ctx context.Context, msg string, attrs ...slog.Attr) {
	cl.LogAttrs(ctx, slog.LevelError, msg, attrs...)
}

// Warn logs a warning message with context attributes
func (cl *ContextLogger) Warn(ctx context.Context, msg string, attrs ...slog.Attr) {
	cl.LogAttrs(ctx, slog.LevelWarn, msg, attrs...)
}

// Debug logs a debug message with context attributes
func (cl *ContextLogger) Debug(ctx context.Context, msg string, attrs ...slog.Attr) {
	cl.LogAttrs(ctx, slog.LevelDebug, msg, attrs...)
}

// With returns a new ContextLogger with additional attributes
func (cl *ContextLogger) With(attrs ...slog.Attr) *ContextLogger {
	// Convert slog.Attr to []any for compatibility with logger.With
	args := make([]any, 0, len(attrs)*2)
	for _, attr := range attrs {
		args = append(args, attr.Key, attr.Value)
	}
	return &ContextLogger{
		logger: cl.logger.With(args...),
	}
}

// WithGroup returns a new ContextLogger with a group name
func (cl *ContextLogger) WithGroup(name string) *ContextLogger {
	return &ContextLogger{
		logger: cl.logger.WithGroup(name),
	}
}

// RequestLogger creates a logger instance configured for a specific request.
// It includes common request attributes like method, path, and remote address.
func RequestLogger(baseLogger *slog.Logger, method, path, remoteAddr string) *slog.Logger {
	return baseLogger.WithGroup("request").With(
		slog.String("method", method),
		slog.String("path", path),
		slog.String("remote_addr", remoteAddr),
	)
}

// ServiceLogger creates a logger instance configured for a specific service.
// It includes service identification attributes.
func ServiceLogger(baseLogger *slog.Logger, serviceName, version string) *slog.Logger {
	return baseLogger.With(
		slog.String("service", serviceName),
		slog.String("version", version),
	)
}
