package humanlog

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// HTTPRequest implements slog.LogValuer for structured HTTP request logging
type HTTPRequest struct {
	Method     string
	URL        string
	RemoteAddr string
	UserAgent  string
	StatusCode int
	Duration   time.Duration
}

// LogValue implements slog.LogValuer interface
func (r HTTPRequest) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("method", r.Method),
		slog.String("url", r.URL),
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("user_agent", r.UserAgent),
		slog.Int("status_code", r.StatusCode),
		slog.Duration("duration", r.Duration),
	)
}

// DatabaseQuery implements slog.LogValuer for structured database query logging
type DatabaseQuery struct {
	Query    string
	Args     []interface{}
	Duration time.Duration
	Error    error
}

// LogValue implements slog.LogValuer interface
func (q DatabaseQuery) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("query", q.Query),
		slog.Duration("duration", q.Duration),
	}

	if len(q.Args) > 0 {
		attrs = append(attrs, slog.Int("arg_count", len(q.Args)))
	}

	if q.Error != nil {
		attrs = append(attrs, slog.String("error", q.Error.Error()))
	}

	return slog.GroupValue(attrs...)
}

// User implements slog.LogValuer for structured user logging
type User struct {
	ID       string
	Username string
	Email    string
	Role     string
}

// LogValue implements slog.LogValuer interface
func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("id", u.ID),
		slog.String("username", u.Username),
		slog.String("email", u.Email),
		slog.String("role", u.Role),
	)
}

// StructuredError provides rich error context for logging
type StructuredError struct {
	Op      string      // operation being performed
	Code    string      // error code
	Message string      // human-readable message
	Cause   error       // underlying error
	Context []slog.Attr // additional context
}

// Error implements the error interface
func (e *StructuredError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Op, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Op, e.Message)
}

// LogValue implements slog.LogValuer interface
func (e *StructuredError) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("op", e.Op),
		slog.String("code", e.Code),
		slog.String("message", e.Message),
	}

	if e.Cause != nil {
		attrs = append(attrs, slog.String("cause", e.Cause.Error()))
	}

	// Add additional context attributes
	attrs = append(attrs, e.Context...)

	return slog.GroupValue(attrs...)
}

// Unwrap returns the underlying error for error wrapping chain
func (e *StructuredError) Unwrap() error {
	return e.Cause
}

// NewStructuredError creates a new structured error with the given operation and message
func NewStructuredError(op, code, message string) *StructuredError {
	return &StructuredError{
		Op:      op,
		Code:    code,
		Message: message,
		Context: make([]slog.Attr, 0),
	}
}

// WithCause adds an underlying cause to the structured error
func (e *StructuredError) WithCause(cause error) *StructuredError {
	e.Cause = cause
	return e
}

// WithContext adds contextual attributes to the error
func (e *StructuredError) WithContext(attrs ...slog.Attr) *StructuredError {
	e.Context = append(e.Context, attrs...)
	return e
}

// LogError logs an error with structured context
func LogError(logger *slog.Logger, err error, msg string, attrs ...slog.Attr) {
	ctx := context.TODO()
	if structErr, ok := err.(*StructuredError); ok {
		// Use the structured error's LogValue for rich context
		allAttrs := append([]slog.Attr{slog.Any("error", structErr)}, attrs...)
		logger.LogAttrs(ctx, slog.LevelError, msg, allAttrs...)
	} else {
		// Fallback to simple error logging
		attrs = append([]slog.Attr{slog.String("error", err.Error())}, attrs...)
		logger.LogAttrs(ctx, slog.LevelError, msg, attrs...)
	}
}

// Example helper functions for common logging patterns

// LogHTTPRequest logs an HTTP request with standardized attributes
func LogHTTPRequest(logger *slog.Logger, req *http.Request, statusCode int, duration time.Duration) {
	httpReq := HTTPRequest{
		Method:     req.Method,
		URL:        req.URL.String(),
		RemoteAddr: req.RemoteAddr,
		UserAgent:  req.UserAgent(),
		StatusCode: statusCode,
		Duration:   duration,
	}

	logger.Info("HTTP request", slog.Any("request", httpReq))
}

// LogDatabaseQuery logs a database query with performance metrics
func LogDatabaseQuery(logger *slog.Logger, query string, args []interface{}, duration time.Duration, err error) {
	dbQuery := DatabaseQuery{
		Query:    query,
		Args:     args,
		Duration: duration,
		Error:    err,
	}

	if err != nil {
		logger.Error("Database query failed", slog.Any("query", dbQuery))
	} else {
		logger.Debug("Database query executed", slog.Any("query", dbQuery))
	}
}

// LogUserAction logs a user action with user context
func LogUserAction(logger *slog.Logger, user User, action, resource string) {
	logger.Info("User action",
		slog.Any("user", user),
		slog.String("action", action),
		slog.String("resource", resource),
	)
}
