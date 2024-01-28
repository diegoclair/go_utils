package logger

import (
	"context"
	"log/slog"
)

type LogParams struct {
	// AppName is the name of your application, that will be used as a field in the log
	AppName string
	// DebugLevel is the level of the log, if true, the log will be in debug level
	DebugLevel bool
	// SlogOptions is the options of the slog library
	SlogOptions slog.HandlerOptions

	// AddAttributesFromContext is a function that will be called to add attributes to the log.
	// it should return a key and a value, example: []any{AccountUUIDKey, ctx.Value(AccountUUIDKey)}
	// example: when you call logger.Info(ctx, "message"), the logger will add the attributes returned by the function
	// it will look like this: {"time":"2020-01-01T00:00:00","level":"INFO","file":"main.go:10","msg":"main: message","account_uuid":"1234567890"}
	AddAttributesFromContext func(ctx context.Context) []any
}

// Logger is a wrapper of the slog library adding some extra functionality
type Logger interface {
	// Info are the same methods as the slog library
	Info(ctx context.Context, msg string)
	// Infof add the functionality of formatting the message, like fmt.Printf
	Infof(ctx context.Context, msg string, args ...any)
	// Infow add the functionality of adding key and values to the log, example: logger.Infow(ctx, "message", "key", "value")
	Infow(ctx context.Context, msg string, keyAndValues ...any)
	// Debug are the same methods as the slog library
	Debug(ctx context.Context, msg string)
	// Debugf add the functionality of formatting the message, like fmt.Printf
	Debugf(ctx context.Context, msg string, args ...any)
	// Debugw add the functionality of adding key and values to the log, example: logger.Debugw(ctx, "message", "key", "value")
	Debugw(ctx context.Context, msg string, keyAndValues ...any)
	// Warn are the same methods as the slog library
	Warn(ctx context.Context, msg string)
	// Warnf add the functionality of formatting the message, like fmt.Printf
	Warnf(ctx context.Context, msg string, args ...any)
	// Warnw add the functionality of adding key and values to the log, example: logger.Warnw(ctx, "message", "key", "value")
	Warnw(ctx context.Context, msg string, keyAndValues ...any)
	// Error are the same methods as the slog library
	Error(ctx context.Context, msg string)
	// Errorf add the functionality of formatting the message, like fmt.Printf
	Errorf(ctx context.Context, msg string, args ...any)
	// Errorw add the functionality of adding key and values to the log, example: logger.Errorw(ctx, "message", "key", "value")
	Errorw(ctx context.Context, msg string, keyAndValues ...any)
	// Fatal are the same methods as the slog library
	Fatal(ctx context.Context, msg string)
	// Fatalw add the functionality of adding key and values to the log, example: logger.Fatalw(ctx, "message", "key", "value")
	Fatalf(ctx context.Context, msg string, args ...any)
	// Fatalw add the functionality of adding key and values to the log, example: logger.Fatalw(ctx, "message", "key", "value")
	Fatalw(ctx context.Context, msg string, keyAndValues ...any)
	// Print implements the SetLog function on mysql library
	Print(args ...any)
	// Printf implements the SetLog function on elasticsearch library
	Printf(msg string, v ...any)
}

// LogParams is the struct that contains the parameters to create a logger
func New(params LogParams) Logger {
	return newSlogLogger(params)
}
