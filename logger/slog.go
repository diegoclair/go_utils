package logger

import (
	"context"
	"fmt"
	"os"

	"log/slog"
)

const (
	LevelFatal     = "FATAL"
	LevelFatalCode = 60
)

var CustomLevels = map[int]string{
	LevelFatalCode: LevelFatal, //high number to avoid conflict with slog levels
}

type SlogLogger struct {
	params LogParams
	*slog.Logger
}

func newSlogLogger(params LogParams) *SlogLogger {
	logger := &SlogLogger{params: params}
	opts := slog.HandlerOptions{}

	if params.DebugLevel {
		opts.Level = slog.LevelDebug
	}

	logger.Logger = slog.New(newCustomJSONFormatter(os.Stdout, params))

	return logger
}

func (l *SlogLogger) Info(ctx context.Context, msg string) {
	l.Logger.InfoContext(ctx, msg, l.params.AddAttributesFromContext(ctx)...)
}

func (l *SlogLogger) Infof(ctx context.Context, msg string, args ...any) {
	l.Logger.InfoContext(ctx, fmt.Sprintf(msg, args...), l.params.AddAttributesFromContext(ctx)...)
}

func (l *SlogLogger) Infow(ctx context.Context, msg string, keyAndValues ...any) {
	l.Logger.InfoContext(ctx, msg, append(l.params.AddAttributesFromContext(ctx), keyAndValues...)...)
}

func (l *SlogLogger) Debug(ctx context.Context, msg string) {
	l.Logger.DebugContext(ctx, msg, l.params.AddAttributesFromContext(ctx)...)
}

func (l *SlogLogger) Debugf(ctx context.Context, msg string, args ...any) {
	l.Logger.DebugContext(ctx, fmt.Sprintf(msg, args...), l.params.AddAttributesFromContext(ctx)...)
}

func (l *SlogLogger) Debugw(ctx context.Context, msg string, keyAndValues ...any) {
	l.Logger.DebugContext(ctx, msg, append(l.params.AddAttributesFromContext(ctx), keyAndValues...)...)
}

func (l *SlogLogger) Warn(ctx context.Context, msg string) {
	l.Logger.WarnContext(ctx, msg, l.params.AddAttributesFromContext(ctx)...)
}

func (l *SlogLogger) Warnf(ctx context.Context, msg string, args ...any) {
	l.Logger.WarnContext(ctx, fmt.Sprintf(msg, args...), l.params.AddAttributesFromContext(ctx)...)
}

func (l *SlogLogger) Warnw(ctx context.Context, msg string, keyAndValues ...any) {
	l.Logger.WarnContext(ctx, msg, append(l.params.AddAttributesFromContext(ctx), keyAndValues...)...)
}

func (l *SlogLogger) Error(ctx context.Context, msg string) {
	l.Logger.ErrorContext(ctx, msg, l.params.AddAttributesFromContext(ctx)...)
}

func (l *SlogLogger) Errorf(ctx context.Context, msg string, args ...any) {
	l.Logger.ErrorContext(ctx, fmt.Sprintf(msg, args...), l.params.AddAttributesFromContext(ctx)...)
}

func (l *SlogLogger) Errorw(ctx context.Context, msg string, keyAndValues ...any) {
	l.Logger.ErrorContext(ctx, msg, append(l.params.AddAttributesFromContext(ctx), keyAndValues...)...)
}

func (l *SlogLogger) Fatal(ctx context.Context, msg string) {
	l.Logger.Log(ctx, LevelFatalCode, msg, l.params.AddAttributesFromContext(ctx)...)
	os.Exit(1)
}

func (l *SlogLogger) Fatalf(ctx context.Context, msg string, args ...any) {
	l.Logger.Log(ctx, LevelFatalCode, fmt.Sprintf(msg, args...), l.params.AddAttributesFromContext(ctx)...)
	os.Exit(1)
}

func (l *SlogLogger) Fatalw(ctx context.Context, msg string, keyAndValues ...any) {
	l.Logger.Log(ctx, LevelFatalCode, msg, append(l.params.AddAttributesFromContext(ctx), keyAndValues...)...)
	os.Exit(1)
}

func (l *SlogLogger) Print(args ...any) {
	l.Logger.Log(context.TODO(), slog.LevelInfo, "", args...)
}

func (l *SlogLogger) Printf(msg string, args ...any) {
	l.Logger.Log(context.TODO(), slog.LevelInfo, fmt.Sprintf(msg, args...))
}
