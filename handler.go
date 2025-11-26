package humanlog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Constants for formatting
const (
	messageWidth = 40 // Fixed width for message field
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
)

// Handler implements slog.Handler for human-readable logging output.
type Handler struct {
	h      slog.Handler
	opts   Options
	mu     sync.Mutex
	attrs  []slog.Attr
	groups []string
}

// Enabled reports whether the handler handles records at the given level.
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

// Handle handles the Record.
// Format: [TIME] LEVEL Message(fixed-width-40-chars) key=value key2=value2 ...
// The message is truncated with ellipsis if it exceeds the fixed width.
// Attributes are displayed in a separate column after the message field.
// If UseJSON is enabled, delegates to the underlying JSON handler.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	// Early return if this level is not enabled (performance optimization)
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	// If JSON mode is enabled, delegate to the underlying handler
	if h.opts.UseJSON {
		return h.h.Handle(ctx, r)
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Format time
	timeStr := r.Time.Format(h.opts.TimeFormat)

	// Format level
	levelStr := formatLevel(r.Level, h.opts.DisableColor)

	// Format message (truncate and pad to configured width)
	message := r.Message
	width := h.opts.MessageWidth
	if width <= 0 {
		width = messageWidth // fallback to constant default
	}
	if len(message) > width {
		// Truncate with ellipsis, ensuring space for "..."
		if width > 3 {
			message = message[:width-3] + "..."
		} else {
			message = message[:width]
		}
	}
	// Use Sprintf with %-*s for left-alignment and padding
	formattedMessage := fmt.Sprintf("%-*s", width, message)

	// Build the log line
	var sb strings.Builder

	// [TIME] LEVEL Message(fixed-width)
	fmt.Fprintf(&sb, "[%s] %s %s", timeStr, levelStr, formattedMessage)

	// Collect and format attributes
	var attrs []string
	attrs = h.appendAttrs(attrs, h.attrs)

	// Add attributes from the record
	r.Attrs(func(attr slog.Attr) bool {
		attrs = h.appendAttrs(attrs, []slog.Attr{attr})
		return true
	})

	// Add source if enabled
	if h.opts.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			shortFile := f.File
			if i := strings.LastIndex(f.File, "/"); i != -1 {
				shortFile = f.File[i+1:]
			}
			attr := slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", shortFile, f.Line))
			attrs = append(attrs, formatAttr(attr, h.opts.DisableColor))
		}
	}

	// Add attributes if any (lazy evaluation - only format if needed)
	if len(attrs) > 0 {
		sb.WriteString(" ")
		sb.WriteString(strings.Join(attrs, " "))
	}

	// Add newline and write to output
	sb.WriteString("\n")
	_, err := io.WriteString(h.opts.Writer, sb.String())
	return err
}

// WithAttrs returns a new Handler whose attributes consist of h's attributes followed by attrs.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if h.opts.UseJSON {
		return &Handler{
			h:      h.h.WithAttrs(attrs),
			opts:   h.opts,
			attrs:  nil,
			groups: nil,
		}
	}

	h2 := *h
	h2.attrs = append(h2.attrs, attrs...)
	return &h2
}

// WithGroup returns a new Handler with the given group name.
func (h *Handler) WithGroup(name string) slog.Handler {
	if h.opts.UseJSON {
		return &Handler{
			h:      h.h.WithGroup(name),
			opts:   h.opts,
			attrs:  nil,
			groups: nil,
		}
	}

	h2 := *h
	h2.groups = append(h2.groups, name)
	return &h2
}

// formatLevel returns a fixed-width, uppercase level string with optional color.
func formatLevel(level slog.Level, disableColor bool) string {
	var levelStr string
	var colorCode string

	switch {
	case level >= slog.LevelError:
		levelStr = "ERROR"
		colorCode = colorRed
	case level >= slog.LevelWarn:
		levelStr = "WARN "
		colorCode = colorYellow
	case level >= slog.LevelInfo:
		levelStr = "INFO "
		colorCode = colorBlue
	default:
		levelStr = "DEBUG"
		colorCode = colorGray
	}

	if disableColor {
		return levelStr
	}

	return fmt.Sprintf("%s%s%s", colorCode, levelStr, colorReset)
}

func (h *Handler) appendAttrs(attrs []string, newAttrs []slog.Attr) []string {
	prefix := strings.Join(h.groups, ".")
	for _, attr := range newAttrs {
		key := attr.Key
		if prefix != "" {
			key = prefix + "." + key
		}
		attrs = append(attrs, formatAttr(slog.Attr{Key: key, Value: attr.Value}, h.opts.DisableColor))
	}
	return attrs
}

// formatAttr formats a single attribute as "key=value".
func formatAttr(attr slog.Attr, disableColor bool) string {
	if attr.Equal(slog.Attr{}) {
		return ""
	}

	key := attr.Key
	val := attr.Value

	// Handle special cases
	switch val.Kind() {
	case slog.KindString:
		// Quote strings if they contain spaces or special characters
		s := val.String()
		if needsQuoting(s) {
			return fmt.Sprintf("%s=%q", key, s)
		}
		return fmt.Sprintf("%s=%s", key, s)

	case slog.KindTime:
		// Format time values
		t := val.Time()
		return fmt.Sprintf("%s=%s", key, t.Format(time.RFC3339))

	case slog.KindDuration:
		// Format duration values
		d := val.Duration()
		return fmt.Sprintf("%s=%s", key, d.String())

	case slog.KindAny:
		// Handle error values specially
		if err, ok := val.Any().(error); ok {
			return fmt.Sprintf("%s=%q", key, err.Error())
		}
		fallthrough

	default:
		// Use the default string representation for other types
		return fmt.Sprintf("%s=%s", key, val.String())
	}
}

// needsQuoting returns true if the string should be quoted in log output.
// Strings are quoted if they:
// - Are empty
// - Contain spaces or control characters
// - Contain special characters that could interfere with log parsing (=, ", ', `, [, ])
// - Look like a Go keyword or boolean (true, false, nil)
// Numbers are not quoted.
func needsQuoting(s string) bool {
	if s == "" {
		return true
	}

	// Don't quote valid numbers
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return false
	}

	// Check for Go keywords and literals that might cause confusion
	switch s {
	case "true", "false", "nil":
		return true
	}

	// Check for spaces, control characters, or special characters
	for _, r := range s {
		if r <= ' ' || r == '=' || r == '"' || r == '\'' || r == '`' || r == '[' || r == ']' || r == '{' || r == '}' {
			return true
		}
	}

	return false
}
