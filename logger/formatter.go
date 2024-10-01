package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gookit/color"
)

type customJSONFormatter struct {
	w         io.Writer
	logToFile bool
	attr      []any
}

func newCustomJSONFormatter(params LogParams) *customJSONFormatter {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error obtaining host name: %v\n", err)
	}

	res := &customJSONFormatter{
		w:         os.Stdout,
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

type orderedEvent struct {
	Time  string                 `json:"time,omitempty"`
	Level string                 `json:"level,omitempty"`
	Msg   string                 `json:"msg,omitempty"`
	File  string                 `json:"file,omitempty"`
	Extra map[string]interface{} `json:"-"`
}

func (oe orderedEvent) MarshalJSON() ([]byte, error) {
	// Start with the fixed fields
	result := fmt.Sprintf(`{"time":"%s","level":"%s","msg":"%s"`, oe.Time, oe.Level, oe.Msg)

	if oe.File != "" {
		result += fmt.Sprintf(`,"file":"%s"`, oe.File)
	}

	// Add extra fields
	for k, v := range oe.Extra {
		jsonValue, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		result += fmt.Sprintf(`,"%s":%s`, k, string(jsonValue))
	}

	result += "}"
	return []byte(result), nil
}

func (f *customJSONFormatter) Write(p []byte) (n int, err error) {
	var rawEvent map[string]interface{}
	if err := json.Unmarshal(p, &rawEvent); err != nil {
		return 0, err
	}

	funcName, fileName, fileLine := f.getRuntimeData()

	event := orderedEvent{
		Time:  getStringValue(rawEvent, "time"),
		Level: strings.ToUpper(getStringValue(rawEvent, "level")),
		Msg:   fmt.Sprintf("%s: %s", funcName, getStringValue(rawEvent, "message")),
		File:  fmt.Sprintf("%s:%d", fileName, fileLine),
		Extra: make(map[string]interface{}),
	}

	if event.Level == "60" {
		event.Level = "CRITICAL"
	}

	// Add custom attributes
	for i := 0; i < len(f.attr); i += 2 {
		if i+1 < len(f.attr) {
			key, ok := f.attr[i].(string)
			if ok {
				event.Extra[key] = f.attr[i+1]
			}
		}
	}

	// Add remaining fields
	for k, v := range rawEvent {
		if k != "time" && k != "level" && k != "message" && k != "caller" {
			event.Extra[k] = v
		}
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		return 0, err
	}

	coloredJSON := f.applyLevelColor(string(jsonData))
	return f.w.Write([]byte(coloredJSON + "\n"))
}

var levelColors = map[string]func(a ...interface{}) string{
	"DEBUG":    color.Magenta.Render,
	"INFO":     color.Blue.Render,
	"WARN":     color.Yellow.Render,
	"ERROR":    color.Red.Render,
	"FATAL":    func(a ...interface{}) string { return color.Bold.Render(color.Red.Render(a...)) },
	"CRITICAL": func(a ...interface{}) string { return color.Bold.Render(color.Red.Render(a...)) },
}

func (f *customJSONFormatter) applyLevelColor(fullMsg string) string {
	if f.logToFile {
		return fullMsg
	}

	for level, colorFunc := range levelColors {
		searchStr := fmt.Sprintf(`"level":"%s"`, level)
		if strings.Contains(fullMsg, searchStr) {
			coloredLevel := colorFunc(level)
			return strings.Replace(fullMsg, searchStr, fmt.Sprintf(`"level":"%s"`, coloredLevel), 1)
		}
	}
	return fullMsg
}

// Helper function to safely get string values from the map
func getStringValue(m map[string]interface{}, key string) string {
	if value, ok := m[key]; ok {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}
	return ""
}

func (f *customJSONFormatter) getRuntimeData() (funcName, filename string, line int) {
	pc, filePath, line, ok := runtime.Caller(8)
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
