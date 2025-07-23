# Humanlog Project Rules

## Shared Guidelines Reference

**This project follows the shared LLM assistant guidelines in [`llm-shared/`](../llm-shared/):**

- **[`llm-shared/project_tech_stack.md`](../llm-shared/project_tech_stack.md)**: Project management, validation, build/test/lint workflows, and universal conventions for all projects.
- **[`llm-shared/languages/go.md`](../llm-shared/languages/go.md)**: Go-specific best practices, library choices, and code formatting/testing standards.
- **[`llm-shared/utils/`](../llm-shared/utils/)**: Tools for code analysis and validation (e.g., `gofuncs`, `validate-docs`).
- **[`llm-shared/templates/`](../llm-shared/templates/)**: Example `.gitignore`, `Taskfile.yml`, and CI workflow templates.

> **Always consult the above files for baseline rules. This file documents project-specific conventions and architectural notes for `humanlog`.**

---

## Project Management & Build Requirements

### Task Completion Criteria

- Task is not complete until `task build` succeeds, which includes:
  - Running tests (`task test`)
  - Linting the code (`task lint`)
  - Building the project (if applicable)
- Task is not complete until it has even basic unit tests, even if they are not comprehensive
  - No need to mock external dependencies, just test the logic of the code
- When working from a markdown checklist of tasks, check off the tasks as you complete them

### Build System Requirements

- Use `task build` over `go build` to ensure all tasks are run
- All build artifacts should be placed in the `build/` directory
- Build tasks must depend on test and lint tasks
- Reference `llm-shared/templates/Taskfile.yml` for comprehensive task structure

### Project Validation

- Use the `validate-docs` tool to check if projects follow standard structure conventions:
  ```bash
  go run llm-shared/utils/validate-docs/validate-docs.go
  ```

---

## Go-Specific Guidelines

### Code Formatting & Quality

- Always run `gofmt -w .` on Go code files after making changes
- Use `go fmt ./...` and `go vet ./...` for linting
- Functions that are easily unit-testable should have tests
- Don't go for 100% test coverage, test the critical parts of the code

### Library Preferences

- Prefer using standard library packages when possible
- Provide justification when adding new third-party dependencies
- Keep dependencies updated
- If SQLite is used, use "modernc.org/sqlite" as the library (no dependency on cgo)
- Logging in applications run from cron: "log/slog" (standard library)
- Logging in applications run from CLI: "fmt.Println" (standard library, use emojis for better UX)
- Configuration management: "github.com/spf13/viper"
- Command-line arguments: "github.com/alecthomas/kong" (only if the project requires complex CLI)

### Function Analysis

- When looking for functions, use the `gofuncs` tool to list all functions in a Go project:
  ```bash
  go run llm-shared/utils/gofuncs/gofuncs.go -dir /path/to/project
  ```

---

## Project-Specific Rules for `humanlog`

### Overview & Architecture

- This package provides a human-readable formatter for Go's `log/slog` output.
- Main components:
  - `Handler` (`handler.go`): Implements `slog.Handler` with custom formatting.
  - `Options` (`options.go`): Configures handler behavior.
  - `NewHandler` (`humanlog.go`): Entry point for handler creation.
- Example usage: `example/main.go`. Tests: `handler_test.go`.

### Developer Workflows

- **Build**: Use `task build` (or standard Go tools `go build`, `go install` if no Taskfile exists yet).
- **Test**: `go test ./...` (see `handler_test.go`).
- **Example**: `go run example/main.go` to see formatted log output.

### Project-Specific Conventions

- No global state; all config via `Options` or method params.
- `NewHandler` panics if given a nil writer (enforced contract).
- All exported types/functions must have doc comments.
- Formatting helpers are unexported and colocated with usage.
- Log messages are padded/truncated to a fixed width (see `messageWidth` in `handler.go`).
- Color output for log levels unless `DisableColor` is set.

### Integration & Dependencies

- Only standard library and `log/slog` are used (following the minimal dependency principle).
- Go module path: `github.com/lepinkainen/humanlog` (see `go.mod`).
- When doing HTTP requests, use a custom user agent that includes the project name and version, e.g. `humanlog/1.0.0`

### Patterns & Examples

- See `example/main.go` for idiomatic usage patterns.
- See `handler_test.go` for test structure and coverage.

---

## Development Workflow

### Code Analysis

- When analyzing large codebases that might exceed context limits, use the Gemini CLI:
  ```bash
  gemini -p "@src/main.go Explain this file's purpose and functionality"
  gemini -p "@src/ Summarise the architecture of this codebase"
  gemini -p "@src/ Is the project test coverage on par with industry standards?"
  ```

### CI/CD Requirements

- Projects should have a basic GitHub Actions setup that uses the build-ci task
- Use `llm-shared/templates/github/workflows/go-ci.yml` as a template
- CI should run tests and linting on push and pull requests
- Use `go test -tags=ci -cover -v ./...` for CI tests
- Allow skipping tests with `//go:build !ci`

### Git Management

- Keep `.gitignore` up to date with Go-specific ignores
- Use `llm-shared/templates/gitignore-go` as a reference
- Ensure build artifacts and temporary files are not committed

---

## Templates & References

### Available Templates

- **Taskfile**: `llm-shared/templates/Taskfile.yml` - Comprehensive task management
- **CI Workflow**: `llm-shared/templates/github/workflows/go-ci.yml` - GitHub Actions for Go
- **Gitignore**: `llm-shared/templates/gitignore-go` - Go-specific ignore patterns
- **Documentation**: `llm-shared/templates/README.md` and `llm-shared/templates/CHANGELOG.md`

### Project Structure

This project follows a simple Go library structure:

```
humanlog/
├── go.mod              # Go module definition
├── *.go                # Main library files (handler.go, options.go, humanlog.go)
├── *_test.go          # Test files
├── example/           # Usage examples
├── llm-shared/        # Shared development guidelines (submodule)
├── docs/              # Project documentation (if needed)
└── build/             # Build artifacts (when using Taskfile)
```

---

_If any section is unclear or missing important project-specific details, please provide feedback or point to additional documentation to improve these rules._
