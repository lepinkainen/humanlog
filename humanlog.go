package humanlog

import (
	"io"
	"log/slog"
)

// NewHandler creates a new human-readable slog.Handler with the given options.
// If opts is nil, default options will be used.
func NewHandler(w io.Writer, opts *Options) *Handler {
	if opts == nil {
		opts = DefaultOptions()
	}

	// Ensure we have a writer
	if w == nil {
		panic("humanlog: nil writer")
	}

	// Set the writer in the options
	options := *opts
	options.Writer = w

	return &Handler{
		h: slog.NewTextHandler(w, &slog.HandlerOptions{
			AddSource:   opts.AddSource,
			Level:       opts.Level,
			ReplaceAttr: nil, // We do our own attribute handling
		}),
		opts:   options,
		attrs:  nil,
		groups: nil,
	}
}

