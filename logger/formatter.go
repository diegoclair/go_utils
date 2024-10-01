package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gookit/color"
	"github.com/rs/zerolog"
)

// customJSONFormatter implementa a interface io.Writer
type customJSONFormatter struct {
	w         io.Writer
	logToFile bool
	attr      []any
}

func newCustomJSONFormatter(w io.Writer, params LogParams) *customJSONFormatter {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error obtaining host name: %v\n", err)
	}

	res := &customJSONFormatter{
		w:         w,
		logToFile: params.LogToFile,
	}

	if params.AppName != "" {
		res.attr = append(res.attr, "app", params.AppName)
	}

	if hostname != "" {
		res.attr = append(res.attr, "host", hostname)
	}

	return res
}

func (f *customJSONFormatter) formatLog(msg string, level zerolog.Level, attributes map[string]interface{}) string {
	funcName, fileName, fileLine := f.getRuntimeData()

	levelStr := strings.ToUpper(level.String())
	if level == CustomLevels[LevelCritical] {
		levelStr = LevelCritical
	}

	// Iniciar com os campos na ordem desejada
	logEntry := fmt.Sprintf(`{"time":"%s","level":"%s","msg":"%s: %s","file":"%s:%d"`,
		time.Now().Format(time.RFC3339),
		levelStr,
		funcName,
		msg,
		fileName,
		fileLine,
	)

	for key, value := range attributes {
		logEntry += fmt.Sprintf(`,"%s":%v`, key, formatValue(value))
	}

	for i := 0; i < len(f.attr); i += 2 {
		if i+1 < len(f.attr) {
			key, ok := f.attr[i].(string)
			if ok {
				logEntry += fmt.Sprintf(`,"%s":%v`, key, formatValue(f.attr[i+1]))
			}
		}
	}

	logEntry += "}"
	return logEntry
}

func (f *customJSONFormatter) Write(p []byte) (n int, err error) {
	coloredJSON := f.applyLevelColor(string(p))
	return f.w.Write(append([]byte(coloredJSON), '\n'))
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

var levelColors = map[string]string{
	"DEBUG":    color.Magenta.Render("DEBUG"),
	"INFO":     color.Blue.Render("INFO"),
	"WARN":     color.Yellow.Render("WARN"),
	"ERROR":    color.Red.Render("ERROR"),
	"FATAL":    color.Bold.Render(color.Red.Render("FATAL")),
	"CRITICAL": color.Bold.Render(color.Red.Render("CRITICAL")),
}

func (f *customJSONFormatter) applyLevelColor(fullMsg string) string {
	if f.logToFile {
		return fullMsg
	}

	for level, colorLevel := range levelColors {
		searchStr := fmt.Sprintf(`"level":"%s"`, level)
		if strings.Contains(fullMsg, searchStr) {
			return strings.Replace(fullMsg, searchStr, fmt.Sprintf(`"level":"%s"`, colorLevel), 1)
		}
	}

	return fullMsg
}

func formatValue(v interface{}) string {
	switch v := v.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf(`"%v"`, v)
	}
}
