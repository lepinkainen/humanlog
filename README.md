# humanlog

A Go package for formatting Go's log/slog output into a more human-readable format.

## Installation

```sh
go get github.com/lepinkainen/humanlog
```

## Usage

```go
import (
    "log/slog"
    "os"
    "github.com/lepinkainen/humanlog"
)

func main() {
    handler := humanlog.NewHandler(os.Stdout, nil) // Use default options
    logger := slog.New(handler)
    logger.Info("Hello, human-readable logs!")
}
```

## Configuration

Pass `nil` as the second argument to use default options:

```go
handler := humanlog.NewHandler(os.Stdout, nil) // Uses all defaults
```

### Default Values

| Option         | Default         | Description                              |
|----------------|-----------------|------------------------------------------|
| `Level`        | `slog.LevelInfo`| Minimum log level                        |
| `TimeFormat`   | `"15:04:05"`    | Timestamp format (HH:MM:SS)              |
| `DisableColor` | `false`         | ANSI color codes enabled                 |
| `AddSource`    | `true`          | Include source file and line number      |
| `MessageWidth` | `40`            | Fixed message width (truncated/padded)   |
| `UseJSON`      | `false`         | Human-readable format (not JSON)         |

### Customizing Options

**Important:** When creating an `Options` struct directly, any unset fields use Go's zero values (empty string, 0, false), not the library defaults. This can cause issues like missing timestamps.

To safely customize options, start from `DefaultOptions()` and modify individual fields:

```go
opts := humanlog.DefaultOptions()
opts.Level = slog.LevelDebug
opts.DisableColor = true
handler := humanlog.NewHandler(os.Stdout, opts)
```

See [`example/main.go`](example/main.go) for more usage patterns.

## Releases

This project follows [Semantic Versioning](https://semver.org/). Releases are automated using GitHub Actions.

- View all releases: [GitHub Releases](https://github.com/lepinkainen/humanlog/releases)
- View changelog: [CHANGELOG.md](CHANGELOG.md)
- Latest version: [![Go Reference](https://pkg.go.dev/badge/github.com/lepinkainen/humanlog.svg)](https://pkg.go.dev/github.com/lepinkainen/humanlog)

### Version History

- `v0.2.x` - Added configurable message width, JSON output, project infrastructure
- `v0.1.0` - Initial release with core functionality
