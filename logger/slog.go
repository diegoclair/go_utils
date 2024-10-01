package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type SlogLogger struct {
	params          LogParams
	logger          zerolog.Logger
	customFormatter *customJSONFormatter
}

func newSlogLogger(params LogParams) *SlogLogger {
	if params.AddAttributesFromContext == nil {
		params.AddAttributesFromContext = func(ctx context.Context) []any {
			return nil
		}
	}

	formatter := newCustomJSONFormatter(os.Stdout, params)
	logger := zerolog.New(formatter).With().Timestamp().Logger()

	if params.DebugLevel {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return &SlogLogger{
		params:          params,
		logger:          logger,
		customFormatter: formatter,
	}
}

func (l *SlogLogger) Close() {
	close(l.customFormatter.done)
}

func (l *SlogLogger) log(ctx context.Context, level zerolog.Level, msg string, fields ...any) {
	attributes := make(map[string]any)

	// add attributes from context
	if l.params.AddAttributesFromContext != nil {
		attrs := l.params.AddAttributesFromContext(ctx)
		for i := 0; i < len(attrs); i += 2 {
			if i+1 < len(attrs) {
				key, ok := attrs[i].(string)
				if ok {
					attributes[key] = attrs[i+1]
				}
			}
		}
	}

	// add additional fields
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				attributes[key] = fields[i+1]
			}
		}
	}

	// send formatted log to channel
	l.customFormatter.logChan <- l.customFormatter.formatLog(msg, level, attributes)
}

func (l *SlogLogger) Info(ctx context.Context, msg string) {
	l.log(ctx, zerolog.InfoLevel, msg)
}

func (l *SlogLogger) Infof(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zerolog.InfoLevel, fmt.Sprintf(msg, args...))
}

func (l *SlogLogger) Infow(ctx context.Context, msg string, keyAndValues ...any) {
	l.log(ctx, zerolog.InfoLevel, msg, keyAndValues...)
}

func (l *SlogLogger) Debug(ctx context.Context, msg string) {
	l.log(ctx, zerolog.DebugLevel, msg)
}

func (l *SlogLogger) Debugf(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zerolog.DebugLevel, fmt.Sprintf(msg, args...))
}

func (l *SlogLogger) Debugw(ctx context.Context, msg string, keyAndValues ...any) {
	l.log(ctx, zerolog.DebugLevel, msg, keyAndValues...)
}

func (l *SlogLogger) Warn(ctx context.Context, msg string) {
	l.log(ctx, zerolog.WarnLevel, msg)
}

func (l *SlogLogger) Warnf(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zerolog.WarnLevel, fmt.Sprintf(msg, args...))
}

func (l *SlogLogger) Warnw(ctx context.Context, msg string, keyAndValues ...any) {
	l.log(ctx, zerolog.WarnLevel, msg, keyAndValues...)
}

func (l *SlogLogger) Error(ctx context.Context, msg string) {
	l.log(ctx, zerolog.ErrorLevel, msg)
}

func (l *SlogLogger) Errorf(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zerolog.ErrorLevel, fmt.Sprintf(msg, args...))
}

func (l *SlogLogger) Errorw(ctx context.Context, msg string, keyAndValues ...any) {
	l.log(ctx, zerolog.ErrorLevel, msg, keyAndValues...)
}

func (l *SlogLogger) Fatal(ctx context.Context, msg string) {
	l.log(ctx, zerolog.FatalLevel, msg)
	os.Exit(1)
}

func (l *SlogLogger) Fatalf(ctx context.Context, msg string, args ...any) {
	l.log(ctx, zerolog.FatalLevel, fmt.Sprintf(msg, args...))
	os.Exit(1)
}

func (l *SlogLogger) Fatalw(ctx context.Context, msg string, keyAndValues ...any) {
	l.log(ctx, zerolog.FatalLevel, msg, keyAndValues...)
	os.Exit(1)
}

func (l *SlogLogger) Critical(ctx context.Context, msg string) {
	l.log(ctx, CustomLevels[LevelCritical], msg)
}

func (l *SlogLogger) Criticalf(ctx context.Context, msg string, args ...any) {
	l.log(ctx, CustomLevels[LevelCritical], fmt.Sprintf(msg, args...))
}

func (l *SlogLogger) Criticalw(ctx context.Context, msg string, keyAndValues ...any) {
	l.log(ctx, CustomLevels[LevelCritical], msg, keyAndValues...)
}

func (l *SlogLogger) Print(args ...any) {
	l.log(context.Background(), zerolog.InfoLevel, fmt.Sprint(args...))
}

func (l *SlogLogger) Printf(msg string, args ...any) {
	l.log(context.Background(), zerolog.InfoLevel, fmt.Sprintf(msg, args...))
}
