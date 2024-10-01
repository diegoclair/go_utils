package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gookit/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
)

type customJSONFormatter struct {
	w         io.Writer
	logToFile bool
	attr      []any
	logChan   chan string
	done      chan struct{}
}

func newCustomJSONFormatter(w io.Writer, params LogParams) *customJSONFormatter {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error obtaining host name: %v\n", err)
	}

	res := &customJSONFormatter{
		w:         w,
		logToFile: false,
		logChan:   make(chan string, 10000),
		done:      make(chan struct{}),
	}

	if params.AppName != "" {
		res.attr = append(res.attr, "app", params.AppName)
	}

	if hostname != "" {
		res.attr = append(res.attr, "host", hostname)
	}

	go res.processLogs()

	return res
}

var builderPool = sync.Pool{
	New: func() any {
		return new(strings.Builder)
	},
}

func (f *customJSONFormatter) formatLog(msg string, level zerolog.Level, attributes map[string]any) string {
	funcName, fileName, fileLine := f.getRuntimeData()

	buf := builderPool.Get().(*strings.Builder)
	buf.Reset()
	defer builderPool.Put(buf)

	buf.WriteString(`{"time":"`)
	buf.WriteString(zerolog.TimestampFunc().Format("2006-01-02T15:04:05"))
	buf.WriteString(`","level":"`)
	buf.WriteString(level.String())
	buf.WriteString(`","file":"`)
	buf.WriteString(fileName)
	buf.WriteString(":")
	buf.WriteString(fmt.Sprintf("%d", fileLine))
	buf.WriteString(`","msg":"`)
	buf.WriteString(funcName)
	buf.WriteString(`: `)
	buf.WriteString(msg)
	buf.WriteString(`"`)

	// Adicionar atributos personalizados
	for i := 0; i < len(f.attr); i += 2 {
		if i+1 < len(f.attr) {
			key, ok := f.attr[i].(string)
			if !ok {
				continue
			}
			value := f.attr[i+1]
			buf.WriteString(`,"`)
			buf.WriteString(key)
			buf.WriteString(`":`)
			buf.WriteString(formatJSONValue(value))
		}
	}

	// Adicionar atributos adicionais
	for key, value := range attributes {
		buf.WriteString(`,"`)
		buf.WriteString(key)
		buf.WriteString(`":`)
		buf.WriteString(formatJSONValue(value))
	}

	buf.WriteString("}")
	return f.applyLevelColor(buf.String(), level.String())
}

func (f *customJSONFormatter) processLogs() {
	for {
		select {
		case logMsg := <-f.logChan:
			_, err := fmt.Fprintln(f.w, logMsg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing log: %v\n", err)
			}
		case <-f.done:
			return
		}
	}
}

func (f *customJSONFormatter) applyLevelColor(fullMsg, level string) string {
	if !f.logToFile {
		levelUpper := strings.ToUpper(level)
		levelColor := ""

		switch level {
		case zerolog.InfoLevel.String():
			levelColor = color.Blue.Render(levelUpper)
		case zerolog.DebugLevel.String():
			levelColor = color.Magenta.Render(levelUpper)
		case zerolog.WarnLevel.String():
			levelColor = color.Yellow.Render(levelUpper)
		case zerolog.ErrorLevel.String():
			levelColor = color.Red.Render(levelUpper)
		case LevelFatal, LevelCritical:
			levelColor = color.Bold.Render(color.Red.Render(levelUpper))
		default:
			levelColor = levelUpper
		}

		return strings.Replace(fullMsg, `"level":"`+level+`"`, `"level":"`+levelColor+`"`, 1)
	}

	return fullMsg
}

func (f *customJSONFormatter) getRuntimeData() (funcName, filename string, line int) {
	pc, filePath, line, ok := runtime.Caller(4)
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

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func formatJSONValue(v any) string {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf(`"%v"`, v)
	}
	return string(jsonBytes)
}

func (f *customJSONFormatter) Write(p []byte) (n int, err error) {
	formattedLog := f.applyLevelColor(string(p), "")
	f.logChan <- formattedLog
	return len(p), nil
}
