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
