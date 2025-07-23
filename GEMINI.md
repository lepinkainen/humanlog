# GEMINI.md for humanlog

This document provides essential guidelines for AI agents working on the `humanlog` Go project.

## Project Overview

`humanlog` is a Go package that provides a human-readable formatter for `log/slog` output.
- **Core Components**:
    - `handler.go`: Implements the custom `slog.Handler` for formatting.
    - `options.go`: Defines configuration options for the handler.
    - `humanlog.go`: Entry point for creating new `humanlog` handlers.
- **Example Usage**: Refer to `example/main.go` for idiomatic usage patterns.

## Developer Workflows

- **Build, Test, Lint**: Always use `task build`. This command orchestrates testing (`go test ./...`), linting, and formatting (`goimports -w .`).
    - **Do NOT** use `go build` or `go test` directly for general tasks, as `task build` ensures all quality checks are performed.
- **Run Example**: To see the formatted log output, execute `go run example/main.go`.
- **Function Analysis**: Use `go run llm-shared/utils/gofuncs/gofuncs.go -dir .` to list functions in the project.

## Project-Specific Conventions

- **Configuration**: All handler configuration is managed via the `Options` struct in `options.go`. No global state is used.
- **Writer Contract**: `humanlog.NewHandler` panics if a `nil` `io.Writer` is provided.
- **Documentation**: All exported types and functions must have doc comments.
- **Formatting**: Log messages are padded/truncated to a fixed width (`messageWidth` in `handler.go`). Color output is enabled by default.
- **Dependencies**: The project prioritizes the Go standard library and `log/slog`. Avoid adding new third-party dependencies without strong justification.

## Integration Points

- `humanlog` integrates directly with Go's standard `log/slog` package by implementing the `slog.Handler` interface.

## General AI Agent Guidance

- **Task Completion**: A task is not complete until `task build` succeeds and basic unit tests are in place.
- **Code Formatting**: Always run `goimports -w .` on Go code files after making changes.
- **Testing**: Focus on testing critical parts of the code; 100% test coverage is not required.
- **CI/CD**: Refer to `llm-shared/templates/github/workflows/go-ci.yml` for CI setup.
- **Git**: Keep `.gitignore` updated using `llm-shared/templates/gitignore-go` as a reference.
