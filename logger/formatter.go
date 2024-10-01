package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"

	"github.com/gookit/color"
	"github.com/labstack/gommon/log"

	"log/slog"
)

type customJSONFormatter struct {
	slog.Handler
	w io.Writer
	// logToFile need to be implemented
	logToFile bool
	attr      []slog.Attr
	logChan   chan string // used to async writing logs
	done      chan struct{}
}

func newCustomJSONFormatter(w io.Writer, params LogParams) *customJSONFormatter {
	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf("error obtaining host name: %v", err)
	}

	res := &customJSONFormatter{
		Handler:   slog.NewJSONHandler(w, &params.slogOptions),
		w:         w,
		logToFile: false,
		logChan:   make(chan string, 1000), // buffered channel to reduce the risk of blocking the main thread
		done:      make(chan struct{}),
	}

	if params.AppName != "" {
		res.attr = append(res.attr, slog.String("app", params.AppName))
	}

	if hostname != "" {
		res.attr = append(res.attr, slog.String("host", hostname))
	}

	// start a new goroutine to write the logs
	go res.processLogs()

	return res

}

// builderPool is a pool of strings.Builder to reduce memory allocation and improve performance
var builderPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

func (f *customJSONFormatter) Handle(ctx context.Context, r slog.Record) error {
	funcName, fileName, fileLine := f.getRuntimeData()

	level := f.getLevel(r.Level)

	buf := builderPool.Get().(*strings.Builder)
	buf.Reset()
	defer builderPool.Put(buf)

	buf.WriteString("{")

	//add default attrs
	buf.WriteString(`"time":"`)
	buf.WriteString(r.Time.Format("2006-01-02T15:04:05"))
	buf.WriteString(`","level":"`)
	buf.WriteString(level)
	buf.WriteString(`","file":"`)
	buf.WriteString(fileName)
	buf.WriteString(":")
	buf.WriteString(fmt.Sprintf("%d", fileLine))
	buf.WriteString(`","msg":"`)
	buf.WriteString(funcName)
	buf.WriteString(`: `)
	buf.WriteString(r.Message)
	buf.WriteString(`"`)

	//add custom attrs
	r.Attrs(func(a slog.Attr) bool {
		buf.WriteString(`,"`)
		buf.WriteString(a.Key)
		buf.WriteString(`":`)
		buf.WriteString(formatJSONValue(a.Value))
		return true
	})

	//add extra attrs
	for _, attr := range f.attr {
		buf.WriteString(`,"`)
		buf.WriteString(attr.Key)
		buf.WriteString(`":`)
		buf.WriteString(attr.Value.String()) // we only have string on f.attr
	}

	buf.WriteString("}")
	f.logChan <- f.applyLevelColor(buf.String(), level)

	return nil
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
			break
		}
	}
}

func (f *customJSONFormatter) WithAttrs(attrs []slog.Attr) slog.Handler {
	return f.Handler.WithAttrs(f.attr)
}

func (f *customJSONFormatter) WithGroup(name string) slog.Handler {
	return f.Handler.WithGroup(name)
}

func (f *customJSONFormatter) Enabled(ctx context.Context, level slog.Level) bool {
	return f.Handler.Enabled(ctx, level)
}

func (f *customJSONFormatter) applyLevelColor(fullMsg, level string) string {

	if !f.logToFile {
		level := level
		levelUpper := strings.ToUpper(level)
		levelColor := ""

		switch level {
		case slog.LevelInfo.String():
			levelColor = color.Blue.Render(levelUpper)
		case slog.LevelDebug.String():
			levelColor = color.Magenta.Render(levelUpper)
		case slog.LevelWarn.String():
			levelColor = color.Yellow.Render(levelUpper)
		case slog.LevelError.String():
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

func (f *customJSONFormatter) getLevel(level slog.Level) string {
	if l, ok := CustomLevels[int(level)]; ok {
		return l
	}
	return level.String()
}

func (f *customJSONFormatter) getRuntimeData() (funcName, filename string, line int) {
	pc, filePath, line, ok := runtime.Caller(5)
	if !ok {
		panic("Could not get context info for logger!")
	}
	filename = filepath.Base(filePath)
	funcPath := runtime.FuncForPC(pc).Name()
	funcName = funcPath[strings.LastIndex(funcPath, ".")+1:]

	//handle go func called inside of a function
	/*
		for example, we have a func Example() and inside of it, we have a go func() without a name, the it will output funcName as func1, with this handle, it will
		output func name as Example.func1
	*/
	if strings.Contains(funcName, "func") {
		funcBefore := funcPath[:strings.LastIndex(funcPath, ".")]
		funcName = funcPath[strings.LastIndex(funcBefore, ".")+1:]
	}
	return
}

// better performance than standard encode/json lib
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func formatJSONValue(v slog.Value) string {
	jsonBytes, err := json.Marshal(v.Any())
	if err != nil {
		// Fallback to string if serialization fails
		return fmt.Sprintf(`"%v"`, v.Any())
	}
	return string(jsonBytes)
}
