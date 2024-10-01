package logger

import (
	"context"
	"fmt"
	"os"

	"log/slog"
)

const (
	// high numbers to avoid conflict with slog levels
	LevelFatal        = "FATAL"
	LevelFatalCode    = 60
	LevelCritical     = "CRITICAL"
	LevelCriticalCode = 61
)

var CustomLevels = map[int]string{
	LevelFatalCode:    LevelFatal,
	LevelCriticalCode: LevelCritical,
}

type SlogLogger struct {
	params LogParams
	*slog.Logger
	customFormatter *customJSONFormatter
}

func newSlogLogger(params LogParams) *SlogLogger {
	if params.AddAttributesFromContext == nil {
		params.AddAttributesFromContext = func(ctx context.Context) []any {
			return nil
		}
	}

	logger := &SlogLogger{params: params, customFormatter: newCustomJSONFormatter(os.Stdout, params)}

	if params.DebugLevel {
		params.slogOptions.Level = slog.LevelDebug
	}

	logger.Logger = slog.New(logger.customFormatter)

	return logger
}

// Close closes the logger, it will flush buffered logs and close log channel
func (l *SlogLogger) Close() {
	close(l.customFormatter.done)
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

func (l *SlogLogger) Critical(ctx context.Context, msg string) {
	l.Logger.Log(ctx, LevelCriticalCode, msg, l.params.AddAttributesFromContext(ctx)...)
}

func (l *SlogLogger) Criticalf(ctx context.Context, msg string, args ...any) {
	l.Logger.Log(ctx, LevelCriticalCode, fmt.Sprintf(msg, args...), l.params.AddAttributesFromContext(ctx)...)
}

func (l *SlogLogger) Criticalw(ctx context.Context, msg string, keyAndValues ...any) {
	l.Logger.Log(ctx, LevelCriticalCode, msg, append(l.params.AddAttributesFromContext(ctx), keyAndValues...)...)
}

func (l *SlogLogger) Print(args ...any) {
	l.Logger.Log(context.TODO(), slog.LevelInfo, "", args...)
}

func (l *SlogLogger) Printf(msg string, args ...any) {
	l.Logger.Log(context.TODO(), slog.LevelInfo, fmt.Sprintf(msg, args...))
}
