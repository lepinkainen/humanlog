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

You can customize log level, color, and time format via the `Options` struct:

```go
opts := &humanlog.Options{Level: slog.LevelDebug, DisableColor: true}
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
