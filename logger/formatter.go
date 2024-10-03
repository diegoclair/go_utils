package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type customJSONFormatter struct {
	logToFile bool
	attr      []zap.Field
}

func newCustomJSONFormatter(params LogParams) *customJSONFormatter {
	res := &customJSONFormatter{
		logToFile: params.LogToFile,
	}

	if params.AppName != "" {
		res.attr = append(res.attr, zap.String("app", params.AppName))
	}

	hostname, err := os.Hostname()
	if err == nil && hostname != "" {
		res.attr = append(res.attr, zap.String("host", hostname))
	}

	return res
}

func (f *customJSONFormatter) Run(entry zapcore.Entry) error {
	funcName, _, _ := f.getRuntimeData()
	entry.Message = fmt.Sprintf("%s: %s", funcName, entry.Message)
	return nil
}

func (f *customJSONFormatter) formatLevel(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	levelStr := l.CapitalString()
	if l == LevelCritical {
		levelStr = "CRITICAL"
	}
	enc.AppendString(levelStr)
}

var levelColors = map[string]func(...any) string{
	"DEBUG":    color.Magenta.Render,
	"INFO":     color.Blue.Render,
	"WARN":     color.Yellow.Render,
	"ERROR":    color.Red.Render,
	"FATAL":    func(a ...any) string { return color.Bold.Render(color.Red.Render(a...)) },
	"CRITICAL": func(a ...any) string { return color.Bold.Render(color.Red.Render(a...)) },
}

func (f *customJSONFormatter) getRuntimeData() (funcName, filename string, line int) {
	pc, filePath, line, ok := runtime.Caller(6)
	if !ok {
		return "unknown", "unknown", 0
	}
	filename = filepath.Base(filePath)
	funcPath := runtime.FuncForPC(pc).Name()
	funcName = funcPath[strings.LastIndex(funcPath, ".")+1:]

	if strings.Contains(funcName, "func") {
		funcBefore := funcPath[:strings.LastIndex(funcPath, ".")]
		funcName = funcPath[strings.LastIndex(funcBefore, ".")+1:]
	}
	return
}

func (f *customJSONFormatter) getDefaultFields() []zap.Field {
	return f.attr
}

func (f *customJSONFormatter) ApplyColors(logEntry string) string {
	if f.logToFile {
		return logEntry
	}

	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(logEntry), &entry); err != nil {
		return logEntry
	}

	level, ok := entry["level"].(string)
	if ok {
		if colorFunc, exists := levelColors[level]; exists {
			entry["level"] = colorFunc(level)
		}
	}

	orderedFields := []string{"time", "level", "msg", "file", "app", "host"}
	var result strings.Builder
	result.Grow(len(logEntry))

	result.WriteString("{")
	firstField := true

	for _, field := range orderedFields {
		if value, exists := entry[field]; exists {
			if !firstField {
				result.WriteString(",")
			}
			firstField = false
			writeField(&result, field, value)
			delete(entry, field)
		}
	}

	for key, value := range entry {
		if !firstField {
			result.WriteString(",")
		}
		firstField = false
		writeField(&result, key, value)
	}

	result.WriteString("}\n")
	return result.String()
}

func writeField(b *strings.Builder, key string, value interface{}) {
	b.WriteString(`"`)
	b.WriteString(key)
	b.WriteString(`":`)
	if key == "level" {
		b.WriteString(`"`)
		b.WriteString(fmt.Sprint(value))
		b.WriteString(`"`)
	} else {
		jsonValue, _ := json.Marshal(value)
		b.Write(jsonValue)
	}
}
