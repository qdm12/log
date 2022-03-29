# Log

*Go simple yet powerful logging library*

[![Build status](https://github.com/qdm12/log/workflows/CI/badge.svg?branch=main)](https://github.com/qdm12/log/actions?query=workflow%3A"CI")
[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/log.svg)](https://github.com/qdm12/log/commits/main)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/log.svg)](https://github.com/qdm12/log/graphs/contributors)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/log.svg)](https://github.com/qdm12/log/issues)

## Focus

This log library focuses on:

- Human readable logs
- Dependency injected loggers
- Thread safety
- Write to multiple writers

## Setup

```sh
go get github.com/qdm12/log
```

## Usage

### Default logger

The logger constructor `log.New()` uses functional options.

By default the logger:

- logs to [`os.Stdout`](https://pkg.go.dev/os#pkg-variables)
- logs at the `INFO` level
- logs the time with the [*RFC3339* format](https://pkg.go.dev/time#pkg-constants)

```go
package main

import "github.com/qdm12/log"

func main() {
    logger := log.New()
    logger.Info("my message")
    // 2022-03-28T10:03:29Z INFO my message
}
```

➡️ [Source code file](examples/default)

### Formatting methods

Each level log method such as `Warn(s string)` has a corresponding formatting method with a trailing `f` such as `Warnf(format string, args ...interface{})`. For example:

```go
package main

import "github.com/qdm12/log"

func main() {
    logger := log.New()
    logger.Warnf("message number %d", 1)
    // 2022-03-29T07:40:12Z WARN message number 1
}
```

➡️ [Source code file](examples/formatting)

### Custom logger

You can customize the logger creation with for example:

```go
package main

import "github.com/qdm12/log"

func main() {
    logger := log.New(
        log.SetLevel(log.LevelDebug),
        log.SetTimeFormat(time.RFC822),
        log.SetWriters(os.Stdout, os.Stderr),
        log.SetComponent("module"),
        log.SetCallerFile(true),
        log.SetCallerFunc(true),
        log.SetCallerLine(true))
    logger.Info("my message")
    // 29 Mar 22 07:16 UTC INFO [module] my message    main.go:L19:main
    // 29 Mar 22 07:16 UTC INFO [module] my message    main.go:L19:main
}
```

➡️ [Source code file](examples/custom)

### Create a logger from a logger

This should be the preferred way to create additional loggers with different settings, since it favors dependency injection.

For example, we create logger `loggerB` from `loggerA` which inherits the settings from `loggerA` and changes the component setting.

```go
package main

import (
    "github.com/qdm12/log"
)

func main() {
    loggerA := log.New(log.SetComponent("A"))
    loggerB := loggerA.New(log.SetComponent("B"))
    loggerA.Info("my message")
    // 2022-03-29T07:35:08Z INFO [A] my message
    loggerB.Info("my message")
    // 2022-03-29T07:35:08Z INFO [B] my message
}
```

Note this is a thread safe operation, and thread safety on the writers is also maintained.

➡️ [Source code file](examples/inherit)

### Create global loggers

You can create multiple loggers with the global constructor `log.New()`, and writers will be thread safe to write to. For example the following won't write to the buffer at the same time:

```go
package main

import (
    "bytes"
    "time"

    "github.com/qdm12/log"
)

func main() {
    writer := bytes.NewBuffer(nil)

    timer := time.NewTimer(time.Second)

    const parallelism = 2
    for i := 0; i < parallelism; i++ {
        go func() {
            logger := log.New(log.SetWriters(writer))
            for {
                select {
                case <-timer.C:
                    return
                default:
                    logger.Info("my message")
                }
            }
        }()
    }
}
```

You can try it with `CGO_ENABLED=1 go run -race ./examples/global` to ensure there is no data race.

➡️ [Source code file](examples/global)

## Detailed features

The following features are available:

- Multiple options available
  - Set the level `DEBUG`, `INFO`, `WARN`, `ERROR`
  - Set time format, for example `time.RFC3339`
  - Set or add one or more `io.Writer`
  - Set a component string
- Create child loggers inheriting configuration
- Thread safe per `io.Writer` for multiple loggers
- Printf-like methods: `Debugf`, `Infof`, `Warnf`, `Errorf`
- Automatic coloring of levels depending on tty
- Safety to use
  - Full unit test coverage
  - End-to-end race tests
  - Used in the following projects:

## Limitations

### Using a lot of writers

This logging library is thread safe for each writer. To achieve, it uses a global map from writer address to mutex pointer.

⚠️ This however means that this map will not shrink at any time.
If you want to use this logging library with **thousands of different writers**, you should not use it as it is, please create an issue.

## Contributing

See [Contributing](.github/CONTRIBUTING.md)

## License

This repository is under an [MIT license](LICENSE)
