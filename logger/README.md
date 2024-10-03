# Logger Package

## Overview

This package provides a powerful and flexible logging wrapper for Go applications, built on top of the `zap` logging library. It offers enhanced functionality and ease of use compared to the standard logging packages.

## Features

- **Colored Log Levels**: Improve log readability with color-coded log levels (e.g., blue for INFO, red for ERROR).
- **Extended Logging Methods**: Includes `f` and `w` variants for all log levels (e.g., `Infof`, `Infow`) for formatted and key-value pair logging.
- **Contextual Logging**: Ability to add default attributes from context using `AddAttributesFromContext(ctx) []LogField`.
- **Customizable JSON Formatting**: Flexible JSON output with custom field ordering and formatting.
- **Multiple Log Levels**: Supports DEBUG, INFO, WARN, ERROR, FATAL, and CRITICAL log levels.
- **File Logging**: Option to log to files in addition to standard output.
- **Performance Optimized**: Utilizes `zap` for high-performance logging.

## Installation
```bash
go get github.com/diegoclair/go_utils
```
## Quick Start

```go
import "github.com/diegoclair/go_utils/logger"

func main() {
    params := logger.LogParams{
        AppName:    "MyApp",
        DebugLevel: true,
        LogToFile:  false,
        AddAttributesFromContext: func(ctx context.Context) []logger.LogField {
            return []logger.LogField{
                logger.String("user_id", getUserIDFromContext(ctx)),
            }
        },
    }
    log := logger.New(params)
    log.Info(ctx, "Application started")
    log.Infow(ctx, "User logged in", logger.Int("user_id", 12345))
}
```


## Usage

### Creating a Logger

Use `logger.New(params)` to create a new logger instance. Configure the logger using `LogParams`:

```go
params := logger.LogParams{
    AppName:    "MyApp",
    DebugLevel: true,
    LogToFile:  false,
    AddAttributesFromContext: func(ctx context.Context) []logger.LogField {
        // Add custom fields from context
    },
}
log := logger.New(params)
```


### Logging Methods

- Basic logging: `log.Info(ctx, "message")`
- Formatted logging: `log.Infof(ctx, "User %s logged in", username)`
- Structured logging: `log.Infow(ctx, "User action", "action", "login", "username", username)`

### Log Levels

- Debug: `log.Debug`, `log.Debugf`, `log.Debugw`
- Info: `log.Info`, `log.Infof`, `log.Infow`
- Warn: `log.Warn`, `log.Warnf`, `log.Warnw`
- Error: `log.Error`, `log.Errorf`, `log.Errorw`
- Fatal: `log.Fatal`, `log.Fatalf`, `log.Fatalw`
- Critical: `log.Critical`, `log.Criticalf`, `log.Criticalw`

### Custom Formatting

The logger uses a custom JSON formatter that allows for colored output and custom field ordering. You can modify the `customJSONFormatter` in the package to adjust the formatting to your needs.

## Advanced Usage

### Contextual Logging

Use the `AddAttributesFromContext` function to automatically add fields from your context to every log entry:

```go
params.AddAttributesFromContext = func(ctx context.Context) []logger.LogField {
    return []logger.LogField{
        logger.String("request_id", getRequestIDFromContext(ctx)),
        logger.String("user_id", getUserIDFromContext(ctx)),
    }
}
```