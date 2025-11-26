package humanlog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"testing"
	"time"
)

// BenchmarkHandler_Handle benchmarks the core Handle method
func BenchmarkHandler_Handle(b *testing.B) {
	h := NewHandler(io.Discard, &Options{
		Level:        slog.LevelInfo,
		DisableColor: true,
		AddSource:    false,
	})
	logger := slog.New(h)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message", slog.Int("count", i), slog.String("status", "running"))
	}
}

// BenchmarkHandler_HandleWithSource benchmarks Handle with source location enabled
func BenchmarkHandler_HandleWithSource(b *testing.B) {
	h := NewHandler(io.Discard, &Options{
		Level:        slog.LevelInfo,
		DisableColor: true,
		AddSource:    true,
	})
	logger := slog.New(h)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message", slog.Int("count", i), slog.String("status", "running"))
	}
}

// BenchmarkHandler_HandleWithColor benchmarks Handle with color enabled
func BenchmarkHandler_HandleWithColor(b *testing.B) {
	h := NewHandler(io.Discard, &Options{
		Level:        slog.LevelInfo,
		DisableColor: false,
		AddSource:    false,
	})
	logger := slog.New(h)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message", slog.Int("count", i), slog.String("status", "running"))
	}
}

// BenchmarkHandler_HandleManyAttrs benchmarks Handle with many attributes
func BenchmarkHandler_HandleManyAttrs(b *testing.B) {
	h := NewHandler(io.Discard, &Options{
		Level:        slog.LevelInfo,
		DisableColor: true,
		AddSource:    false,
	})
	logger := slog.New(h)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message",
			slog.Int("count", i),
			slog.String("status", "running"),
			slog.Duration("elapsed", time.Millisecond*100),
			slog.Bool("success", true),
			slog.String("user", "test-user"),
			slog.Int64("timestamp", time.Now().Unix()),
		)
	}
}

// BenchmarkHandler_HandleWithGroup benchmarks Handle with grouped attributes
func BenchmarkHandler_HandleWithGroup(b *testing.B) {
	h := NewHandler(io.Discard, &Options{
		Level:        slog.LevelInfo,
		DisableColor: true,
		AddSource:    false,
	})
	groupedHandler := h.WithGroup("request")
	logger := slog.New(groupedHandler)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message",
			slog.String("method", "GET"),
			slog.String("path", "/api/v1/users"),
			slog.Int("status", 200),
		)
	}
}

// BenchmarkHandler_HandleWithAttrs benchmarks Handle with pre-registered attributes
func BenchmarkHandler_HandleWithAttrs(b *testing.B) {
	h := NewHandler(io.Discard, &Options{
		Level:        slog.LevelInfo,
		DisableColor: true,
		AddSource:    false,
	})
	attrsHandler := h.WithAttrs([]slog.Attr{
		slog.String("service", "api"),
		slog.String("version", "1.0.0"),
		slog.String("env", "production"),
	})
	logger := slog.New(attrsHandler)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info("Benchmark message", slog.Int("request_id", i))
	}
}

// BenchmarkHandler_Enabled benchmarks the Enabled method
func BenchmarkHandler_Enabled(b *testing.B) {
	h := NewHandler(io.Discard, &Options{
		Level:        slog.LevelInfo,
		DisableColor: true,
		AddSource:    false,
	})
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		h.Enabled(ctx, slog.LevelInfo)
		h.Enabled(ctx, slog.LevelDebug)
		h.Enabled(ctx, slog.LevelWarn)
		h.Enabled(ctx, slog.LevelError)
	}
}

// BenchmarkFormatAttr benchmarks the formatAttr function with different attribute types
func BenchmarkFormatAttr(b *testing.B) {
	testCases := []struct {
		name string
		attr slog.Attr
	}{
		{"String", slog.String("key", "value")},
		{"Int", slog.Int("key", 42)},
		{"Bool", slog.Bool("key", true)},
		{"Time", slog.Time("key", time.Now())},
		{"Duration", slog.Duration("key", time.Second)},
		{"QuotedString", slog.String("key", "value with spaces")},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				formatAttr(tc.attr, true)
			}
		})
	}
}

// BenchmarkNeedsQuoting benchmarks the needsQuoting function
func BenchmarkNeedsQuoting(b *testing.B) {
	testStrings := []string{
		"simple",
		"with spaces",
		"true",
		"false",
		"nil",
		"123",
		"value=with=equals",
		"",
	}

	for _, s := range testStrings {
		b.Run(s, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				needsQuoting(s)
			}
		})
	}
}

// BenchmarkHandler_ConfigurableWidth benchmarks different message widths
func BenchmarkHandler_ConfigurableWidth(b *testing.B) {
	widths := []int{20, 40, 80, 120}

	for _, width := range widths {
		b.Run(fmt.Sprintf("Width_%d", width), func(b *testing.B) {
			h := NewHandler(io.Discard, &Options{
				Level:        slog.LevelInfo,
				MessageWidth: width,
				DisableColor: true,
				AddSource:    false,
			})
			logger := slog.New(h)

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				logger.Info("This is a test message that may be truncated depending on the configured width",
					slog.Int("iteration", i))
			}
		})
	}
}

// BenchmarkHandler_JSONvs Human benchmarks JSON vs human-readable formats
func BenchmarkHandler_JSONvsHuman(b *testing.B) {
	b.Run("JSON", func(b *testing.B) {
		h := NewHandler(io.Discard, &Options{
			Level:        slog.LevelInfo,
			UseJSON:      true,
			DisableColor: true,
			AddSource:    false,
		})
		logger := slog.New(h)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			logger.Info("Benchmark message",
				slog.String("format", "json"),
				slog.Int("iteration", i),
				slog.Duration("elapsed", time.Millisecond*100),
			)
		}
	})

	b.Run("Human", func(b *testing.B) {
		h := NewHandler(io.Discard, &Options{
			Level:        slog.LevelInfo,
			UseJSON:      false,
			DisableColor: true,
			AddSource:    false,
		})
		logger := slog.New(h)

		b.ResetTimer()
		b.ReportAllocs()

		for i := 0; i < b.N; i++ {
			logger.Info("Benchmark message",
				slog.String("format", "human"),
				slog.Int("iteration", i),
				slog.Duration("elapsed", time.Millisecond*100),
			)
		}
	})
}

// BenchmarkMiddleware_ContextLogger benchmarks the context logger middleware
func BenchmarkMiddleware_ContextLogger(b *testing.B) {
	h := NewHandler(io.Discard, &Options{
		Level:        slog.LevelInfo,
		DisableColor: true,
		AddSource:    false,
	})
	baseLogger := slog.New(h)
	contextLogger := NewContextLogger(baseLogger)

	ctx := WithRequestID(context.Background(), "req-123")
	ctx = WithTraceID(ctx, "trace-456")
	ctx = WithUserID(ctx, "user-789")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		contextLogger.Info(ctx, "Context logger benchmark", slog.Int("iteration", i))
	}
}

// BenchmarkLogValuer benchmarks custom LogValuer implementations
func BenchmarkLogValuer(b *testing.B) {
	h := NewHandler(io.Discard, &Options{
		Level:        slog.LevelInfo,
		DisableColor: true,
		AddSource:    false,
	})
	logger := slog.New(h)

	httpReq := HTTPRequest{
		Method:     "GET",
		URL:        "/api/v1/users/123",
		RemoteAddr: "192.168.1.1:12345",
		UserAgent:  "Mozilla/5.0 Test Agent",
		StatusCode: 200,
		Duration:   time.Millisecond * 150,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		logger.Info("HTTP request processed", slog.Any("request", httpReq))
	}
}
