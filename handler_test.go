package humanlog

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestHandler_Enabled(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
		want  bool
	}{
		{
			name:  "Debug level with Info handler",
			level: slog.LevelDebug,
			want:  false,
		},
		{
			name:  "Info level with Info handler",
			level: slog.LevelInfo,
			want:  true,
		},
		{
			name:  "Warn level with Info handler",
			level: slog.LevelWarn,
			want:  true,
		},
		{
			name:  "Error level with Info handler",
			level: slog.LevelError,
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(new(bytes.Buffer), &Options{Level: slog.LevelInfo})
			if got := h.Enabled(context.Background(), tt.level); got != tt.want {
				t.Errorf("Handler.Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name       string
		level      slog.Level
		message    string
		attrs      []slog.Attr
		wantPrefix string
		wantAttrs  []string
	}{
		{
			name:       "Info message without attributes",
			level:      slog.LevelInfo,
			message:    "Test message",
			attrs:      nil,
			wantPrefix: "INFO",
			wantAttrs:  nil,
		},
		{
			name:       "Error message with string attribute",
			level:      slog.LevelError,
			message:    "Error occurred",
			attrs:      []slog.Attr{slog.String("error", "file not found")},
			wantPrefix: "ERROR",
			wantAttrs:  []string{`error="file not found"`},
		},
		{
			name:       "Warn message with multiple attributes",
			level:      slog.LevelWarn,
			message:    "Warning",
			attrs:      []slog.Attr{slog.Int("count", 42), slog.Bool("critical", false)},
			wantPrefix: "WARN",
			wantAttrs:  []string{"count=42", "critical=false"},
		},
		{
			name:       "Debug message with complex attributes",
			level:      slog.LevelDebug,
			message:    "Debug info",
			attrs:      []slog.Attr{slog.Time("timestamp", time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)), slog.Duration("elapsed", 5*time.Second)},
			wantPrefix: "DEBUG",
			wantAttrs:  []string{"timestamp=", "elapsed=5s"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			h := NewHandler(buf, &Options{Level: slog.LevelDebug, DisableColor: true})
			logger := slog.New(h)

			if len(tt.attrs) > 0 {
				logger.LogAttrs(context.Background(), tt.level, tt.message, tt.attrs...)
			} else {
				logger.Log(context.Background(), tt.level, tt.message)
			}

			got := buf.String()

			// Check level prefix
			if !strings.Contains(got, tt.wantPrefix) {
				t.Errorf("Handler.Handle() output = %v, should contain level %v", got, tt.wantPrefix)
			}

			// Check message
			if !strings.Contains(got, tt.message) {
				t.Errorf("Handler.Handle() output = %v, should contain message %v", got, tt.message)
			}

			// Check attributes if any
			for _, attr := range tt.wantAttrs {
				if !strings.Contains(got, attr) {
					t.Errorf("Handler.Handle() output = %v, should contain attribute %v", got, attr)
				}
			}

			// Check that the message is padded to the fixed width
			// This is a simple check to ensure the fixed-width formatting is applied
			defaultWidth := 40 // default message width
			if len(tt.message) < defaultWidth {
				// If the message is shorter than defaultWidth, there should be spaces after it
				expectedPadding := defaultWidth - len(tt.message)
				if !strings.Contains(got, tt.message+strings.Repeat(" ", expectedPadding)) {
					t.Errorf("Handler.Handle() output = %v, should contain padded message", got)
				}
			}
		})
	}
}

func TestHandler_WithAttrs(t *testing.T) {
	buf := new(bytes.Buffer)
	h := NewHandler(buf, &Options{Level: slog.LevelInfo, DisableColor: true})

	// Create a handler with pre-defined attributes
	h2 := h.WithAttrs([]slog.Attr{slog.String("component", "test"), slog.Int("id", 123)})
	logger := slog.New(h2)

	// Log a message
	logger.Info("Test with attrs")

	got := buf.String()

	// Check that pre-defined attributes are included
	if !strings.Contains(got, "component=test") {
		t.Errorf("WithAttrs() output = %v, should contain component=test", got)
	}

	if !strings.Contains(got, "id=123") {
		t.Errorf("WithAttrs() output = %v, should contain id=123", got)
	}
}

func TestHandler_WithGroup(t *testing.T) {
	buf := new(bytes.Buffer)
	h := NewHandler(buf, &Options{Level: slog.LevelInfo, DisableColor: true})

	// Create a handler with a group
	h2 := h.WithGroup("request")

	// Add attributes to the group
	h3 := h2.WithAttrs([]slog.Attr{slog.String("method", "GET"), slog.String("path", "/api/v1/users")})

	logger := slog.New(h3)

	// Log a message
	logger.Info("Received request")

	got := buf.String()

	// Check that grouped attributes are included with the correct prefix
	if !strings.Contains(got, "request.method=GET") {
		t.Errorf("WithGroup() output = %v, should contain request.method=GET", got)
	}

	// The path contains a slash, which may or may not be quoted depending on the implementation
	if !strings.Contains(got, "request.path=/api/v1/users") && !strings.Contains(got, "request.path=\"/api/v1/users\"") {
		t.Errorf("WithGroup() output = %v, should contain request.path=/api/v1/users or request.path=\"/api/v1/users\"", got)
	}
}

func TestHandler_CustomTimeFormat(t *testing.T) {
	buf := new(bytes.Buffer)
	h := NewHandler(buf, &Options{
		Level:        slog.LevelInfo,
		TimeFormat:   "2006-01-02",
		DisableColor: true,
	})

	logger := slog.New(h)
	logger.Info("Custom time format test")

	got := buf.String()

	// Get today's date in the format YYYY-MM-DD
	today := time.Now().Format("2006-01-02")

	// Check that the log contains today's date
	if !strings.Contains(got, today) {
		t.Errorf("Custom time format output = %v, should contain date %v", got, today)
	}
}

func TestHandler_AttributeFormatting(t *testing.T) {
	tests := []struct {
		name     string
		attr     slog.Attr
		expected string
	}{
		{
			name:     "Empty string",
			attr:     slog.String("key", ""),
			expected: `key=""`,
		},
		{
			name:     "String with spaces",
			attr:     slog.String("key", "hello world"),
			expected: `key="hello world"`,
		},
		{
			name:     "String without spaces",
			attr:     slog.String("key", "hello"),
			expected: `key=hello`,
		},
		{
			name:     "Integer",
			attr:     slog.Int("key", 42),
			expected: `key=42`,
		},
		{
			name:     "Boolean true",
			attr:     slog.Bool("key", true),
			expected: `key=true`,
		},
		{
			name:     "Boolean false",
			attr:     slog.Bool("key", false),
			expected: `key=false`,
		},
		{
			name:     "Error",
			attr:     slog.Any("key", errors.New("test error")),
			expected: `key="test error"`,
		},
		{
			name:     "String with special characters",
			attr:     slog.String("key", "value{with}braces"),
			expected: `key="value{with}braces"`,
		},
		{
			name:     "String that looks like Go keyword",
			attr:     slog.String("key", "nil"),
			expected: `key="nil"`,
		},
		{
			name:     "String that looks like boolean literal",
			attr:     slog.String("key", "true"),
			expected: `key="true"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			h := NewHandler(buf, &Options{Level: slog.LevelInfo, DisableColor: true})
			logger := slog.New(h)

			logger.LogAttrs(context.Background(), slog.LevelInfo, "Test", tt.attr)

			got := buf.String()

			if !strings.Contains(got, tt.expected) {
				t.Errorf("Attribute formatting output = %v, should contain %v", got, tt.expected)
			}
		})
	}
}

func TestHandler_ConfigurableMessageWidth(t *testing.T) {
	tests := []struct {
		name         string
		messageWidth int
		message      string
		expectedLen  int
	}{
		{
			name:         "Custom width 20",
			messageWidth: 20,
			message:      "Short message",
			expectedLen:  20,
		},
		{
			name:         "Custom width 60",
			messageWidth: 60,
			message:      "This is a longer message that should be padded",
			expectedLen:  60,
		},
		{
			name:         "Message too long for width 10",
			messageWidth: 10,
			message:      "This message is way too long for the width",
			expectedLen:  10,
		},
		{
			name:         "Zero width fallback to default",
			messageWidth: 0,
			message:      "Test message",
			expectedLen:  40, // should fallback to default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			h := NewHandler(buf, &Options{
				Level:        slog.LevelInfo,
				MessageWidth: tt.messageWidth,
				DisableColor: true,
			})
			logger := slog.New(h)

			logger.Info(tt.message)
			got := buf.String()

			// Extract the message part from the log line
			// Format: [TIME] LEVEL MESSAGE(fixed-width) attributes...
			parts := strings.Split(got, "] ")
			if len(parts) < 2 {
				t.Errorf("Unexpected log format: %v", got)
				return
			}

			// Get the part after "] " which contains: LEVEL MESSAGE attributes...
			levelMessagePart := parts[1]

			// Split by spaces to separate level from message
			fields := strings.Fields(levelMessagePart)
			if len(fields) < 1 {
				t.Errorf("Unexpected format after timestamp: %v", got)
				return
			}

			// Find where the message starts (after level)
			levelEndIndex := strings.Index(levelMessagePart, fields[0]) + len(fields[0])
			if levelEndIndex >= len(levelMessagePart) {
				t.Errorf("Could not find message part: %v", got)
				return
			}

			// Skip spaces after level
			messageStart := levelEndIndex
			for messageStart < len(levelMessagePart) && levelMessagePart[messageStart] == ' ' {
				messageStart++
			}

			// Extract the fixed-width message (it should be exactly tt.expectedLen characters)
			if messageStart+tt.expectedLen > len(levelMessagePart) {
				t.Errorf("Message part shorter than expected: %v", got)
				return
			}

			messagePart := levelMessagePart[messageStart : messageStart+tt.expectedLen]

			// The messagePart should be exactly the expected length
			if len(messagePart) != tt.expectedLen {
				t.Errorf("Message width = %d, expected %d. Message: %q", len(messagePart), tt.expectedLen, messagePart)
			}

			// Check truncation with ellipsis for long messages
			if len(tt.message) > tt.messageWidth && tt.messageWidth > 3 {
				if !strings.HasSuffix(strings.TrimSpace(messagePart), "...") {
					t.Errorf("Long message should be truncated with ellipsis: %q", messagePart)
				}
			}
		})
	}
}

func TestHandler_JSONOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	h := NewHandler(buf, &Options{
		Level:        slog.LevelInfo,
		UseJSON:      true,
		DisableColor: true,
	})
	logger := slog.New(h)

	logger.Info("Test JSON message", slog.String("key", "value"), slog.Int("count", 42))
	got := buf.String()

	// JSON output should contain JSON structure
	if !strings.Contains(got, `"msg":"Test JSON message"`) {
		t.Errorf("JSON output should contain message field: %v", got)
	}
	if !strings.Contains(got, `"key":"value"`) {
		t.Errorf("JSON output should contain attributes: %v", got)
	}
	if !strings.Contains(got, `"count":42`) {
		t.Errorf("JSON output should contain integer attribute: %v", got)
	}

	// JSON output should contain level field instead of human-readable format markers
	if strings.Contains(got, `"level":"INFO"`) {
		// This is correct - JSON format uses "level":"INFO"
	} else if strings.Contains(got, "[") && strings.Contains(got, "] INFO") {
		t.Errorf("JSON output should not contain human-readable format: %v", got)
	}
}

func TestHandler_JSONWithGroupsAndAttrs(t *testing.T) {
	buf := new(bytes.Buffer)
	h := NewHandler(buf, &Options{
		Level:   slog.LevelInfo,
		UseJSON: true,
	})

	// Test grouped attributes in JSON mode
	groupedHandler := h.WithGroup("request").WithAttrs([]slog.Attr{
		slog.String("id", "123"),
	})
	logger := slog.New(groupedHandler)

	logger.Info("Request processed", slog.String("status", "success"))
	got := buf.String()

	// Should contain grouped attributes in JSON format
	if !strings.Contains(got, `"request"`) {
		t.Errorf("JSON output should contain group: %v", got)
	}
}

func TestHandler_ColorOutput(t *testing.T) {
	// Test with color enabled
	bufWithColor := new(bytes.Buffer)
	hWithColor := NewHandler(bufWithColor, &Options{
		Level:        slog.LevelInfo,
		DisableColor: false,
	})

	// Test with color disabled
	bufNoColor := new(bytes.Buffer)
	hNoColor := NewHandler(bufNoColor, &Options{
		Level:        slog.LevelInfo,
		DisableColor: true,
	})

	// Log the same message with both handlers
	loggerWithColor := slog.New(hWithColor)
	loggerNoColor := slog.New(hNoColor)

	loggerWithColor.Error("Test error message")
	loggerNoColor.Error("Test error message")

	gotWithColor := bufWithColor.String()
	gotNoColor := bufNoColor.String()

	// The colored output should contain ANSI escape codes
	if !strings.Contains(gotWithColor, colorRed) {
		t.Errorf("Color output = %v, should contain ANSI color code", gotWithColor)
	}

	// The non-colored output should not contain ANSI escape codes
	if strings.Contains(gotNoColor, colorRed) {
		t.Errorf("Non-color output = %v, should not contain ANSI color code", gotNoColor)
	}
}
