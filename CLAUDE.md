# CLAUDE.md for humanlog

Essential guide for AI coding agents working on the `humanlog` Go project.

## Project Overview

`humanlog` is a Go package that provides a human-readable formatter for Go's `log/slog` output. It implements the `slog.Handler` interface to format logs with fixed-width messages, colored levels, and structured attributes.

**Core Architecture:**
- `handler.go`: Main `Handler` type implementing `slog.Handler` with custom formatting (234 lines)
- `options.go`: `Options` struct for handler configuration 
- `humanlog.go`: Entry point `NewHandler()` function with nil writer panic contract
- `example/main.go`: Comprehensive usage patterns and API demonstration

**Key Design Patterns:**
- Fixed-width message formatting (40 chars) with truncation/padding in `handler.go:57-64`
- No global state - all configuration via `Options` struct or method parameters
- Immutable handler creation via `WithAttrs()` and `WithGroup()` methods
- Structured attribute grouping with dot notation (e.g., `request.method=GET`)

## Developer Workflows

**Build & Test Commands:**
- **Primary**: `task build` - runs tests, linting, and formatting (`goimports -w .`)
- **Fallback**: Standard Go tools (`go test ./...`, `go build`) if no Taskfile exists
- **Example**: `go run example/main.go` to see formatted output
- **Analysis**: `go run llm-shared/utils/gofuncs/gofuncs.go -dir .` for function listing

**Task Completion Criteria:**
- Must run `gofmt -w` on changed Go files before build attempts
- Task incomplete until `task build` succeeds
- Basic unit tests required (see `handler_test.go` for patterns)

## Project-Specific Conventions

**Handler Contract:**
- `NewHandler(nil, opts)` panics - enforced in `humanlog.go:16-18`
- All exported functions require doc comments
- Message width fixed at 40 characters with ellipsis truncation
- Color output enabled by default unless `DisableColor: true`

**Attribute Formatting:**
- Strings with spaces/special chars are quoted: `key="value with spaces"`
- Time values use RFC3339 format
- Error values are quoted: `error="connection refused"`
- Go keywords (`true`, `false`, `nil`) are quoted when used as string values

**Testing Patterns:**
- Use `bytes.Buffer` for output capture in tests
- Set `DisableColor: true` for predictable test assertions
- Test both grouped and ungrouped attribute scenarios
- Cover time formatting and attribute quoting edge cases

## Integration Points

**slog Integration:**
- Implements `slog.Handler` interface (Enable, Handle, WithAttrs, WithGroup)
- Uses embedded `slog.TextHandler` for level filtering and source location
- Compatible with `slog.SetDefault()` for global logger replacement

**Dependencies:**
- Standard library only (no external dependencies in `go.mod`)
- Go 1.24.5+ required
- Module path: `github.com/lepinkainen/humanlog`

## Development Guidelines

**Code Quality:**
- Follow `llm-shared/` conventions for build, lint, and test standards
- Prefer standard library over third-party dependencies
- Use `llm-shared/utils/validate-docs/` to verify project structure
- Reference `llm-shared/templates/` for CI, gitignore, and build templates

**Formatting Examples:**
```
[15:04:05] INFO  User logged in successfully        user_id=123 session="abc-def" request.ip=192.168.1.1
[15:04:05] ERROR Connection failed                  error="timeout after 30s" retries=3 source=main.go:42
```