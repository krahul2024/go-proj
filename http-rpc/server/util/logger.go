package util

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

const (
    Reset  = "\033[0m"
    Red    = "\033[31m"
    Green  = "\033[32m"
    Yellow = "\033[33m"
    Blue   = "\033[34m"
    Purple = "\033[35m"
    Cyan   = "\033[36m"
    Gray   = "\033[37m"
    White  = "\033[97m"
)

type Logger struct {
    logger zerolog.Logger
    cwd    string
}

var globalLogger *Logger

func init() {
    globalLogger = New()
}

func New() *Logger {
    consoleWriter := &CustomConsoleWriter{
        Out:        os.Stdout,
        TimeFormat: time.RFC3339,
    }

    logger   := zerolog.New(consoleWriter).With().Timestamp().Logger()
    cwd, err := os.Getwd()

    if err != nil {
        log.Fatalln("Error init'ng logger", err)
        os.Exit(1)
    }

    return &Logger{logger: logger, cwd: cwd}
}

type CustomConsoleWriter struct {
    Out        *os.File
    TimeFormat string
}

func (w *CustomConsoleWriter) Write(p []byte) (n int, err error) {
    var event map[string]interface{}
    if err := json.Unmarshal(p, &event); err != nil {
        return w.Out.Write(p)
    }

    pid := os.Getpid()

    timeStr := ""
    if ts, ok := event["time"].(string); ok {
        if t, err := time.Parse(time.RFC3339, ts); err == nil {
            timeStr = t.UTC().Format("2006-01-02 15:04:05")
        }
    }

    level := "INFO"
    levelColor := White
    if l, ok := event["level"].(string); ok {
        level = strings.ToUpper(l)
        switch level {
        case "TRACE":
            levelColor = Gray
        case "DEBUG":
            levelColor = Blue
        case "INFO":
            levelColor = Green
        case "WARN":
            levelColor = Yellow
        case "ERROR":
            levelColor = Red
        case "FATAL":
            levelColor = Purple
        case "PANIC":
            levelColor = Purple
        }
    }

    message := ""
    if msg, ok := event["message"].(string); ok {
        message = msg
    }

    caller := ""
    if c, ok := event["caller"].(string); ok {
        caller = c
    }

    var kvPairs []string
    for k, v := range event {
        if k != "time" && k != "level" && k != "message" && k != "caller" {
            kvPairs = append(kvPairs, fmt.Sprintf("%s = %v", k, v))
        }
    }

    output := fmt.Sprintf("[%d] %s %s[%s]%s %s",
        pid,
        timeStr,
        levelColor,
        level,
        Reset,
        message,
        )

    if len(kvPairs) > 0 {
        output += fmt.Sprintf(" (%s)", strings.Join(kvPairs, ", "))
    }

    if caller != "" {
        output += fmt.Sprintf(" %s", caller)
    }

    output += "\n"

    return w.Out.Write([]byte(output))
}

func getCaller(cwd string) string {
    _, file, line, ok := runtime.Caller(4)
    if !ok {
        return ""
    }

    relFilePath, err := filepath.Rel(cwd, file)
    if err != nil {
        log.Fatalln("Error getting current dir", err)
    }

    return fmt.Sprintf("%s:%d", relFilePath, line)
}

func (l *Logger) logWithFields(level zerolog.Level, msg string, fields map[string]interface{}) {
    event := l.logger.WithLevel(level).Str("caller", getCaller(l.cwd))

    for k, v := range fields {
        event = event.Interface(k, v)
    }

    event.Msg(msg)
}

func (l *Logger) Trace(msg string, fields ...map[string]interface{}) {
    kvFields := flatten(fields)
    l.logWithFields(zerolog.TraceLevel, msg, kvFields)
}

func (l *Logger) Debug(msg string, fields ...map[string]interface{}) {
    kvFields := flatten(fields)
    l.logWithFields(zerolog.DebugLevel, msg, kvFields)
}

func (l *Logger) Info(msg string, fields ...map[string]interface{}) {
    kvFields := flatten(fields)
    l.logWithFields(zerolog.InfoLevel, msg, kvFields)
}

func (l *Logger) Warn(msg string, fields ...map[string]interface{}) {
    kvFields := flatten(fields)
    l.logWithFields(zerolog.WarnLevel, msg, kvFields)

    if kvFields["stack"] == true {
        fmt.Println(Yellow, string(debug.Stack()), Reset)
    }
}

func (l *Logger) Error(err error, msg string, fields ...map[string]interface{}) {
    kvFields               := flatten(fields)
    kvFields["error"]       = err.Error()
    l.logWithFields(zerolog.ErrorLevel, msg, kvFields)
    fmt.Println(Red, string(debug.Stack()), Reset)
}

func (l *Logger) ErrorMsg(msg string, fields ...map[string]interface{}) {
    l.logWithFields(zerolog.ErrorLevel, msg, flatten(fields))
}

func (l *Logger) Fatal(err error, msg string, fields ...map[string]interface{}) {
    kvFields := flatten(fields)
    l.logWithFields(zerolog.FatalLevel, msg, kvFields)
    os.Exit(1)
}

func (l *Logger) Panic(err error, msg string, fields ...map[string]interface{}) {
    kvFields := flatten(fields)
    l.logWithFields(zerolog.PanicLevel, msg, kvFields)
    panic(msg)
}

func Trace(msg string, fields ...map[string]interface{}) {
    globalLogger.Trace(msg, fields...)
}

func Debug(msg string, fields ...map[string]interface{}) {
    globalLogger.Debug(msg, fields...)
}

func ErrMsg(msg string, fields ...map[string]interface{}) {
    globalLogger.Error(nil, msg, fields...)
}

func Error(err error, msg string, fields ...map[string]interface{}) {
    globalLogger.Error(err, msg, fields...)
}

func Info(msg string, fields ...map[string]interface{}) {
    globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...map[string]interface{}) {
    globalLogger.Warn(msg, fields...)
}

func Fatal(err error, msg string, fields ...map[string]interface{}) {
    globalLogger.Fatal(err, msg, fields...)
}

func Panic(err error, msg string, fields ...map[string]interface{}) {
    globalLogger.Panic(err, msg, fields...)
}

func SetLevel(level zerolog.Level) {
    globalLogger.logger = globalLogger.logger.Level(level)
}

func GetLogger() *Logger {
    return globalLogger
}

func flatten(fields []map[string]interface{}) map[string]interface{} {
    merged := make(map[string]interface{})

    for _, field := range fields {
        for k, v := range field {
            merged[k] = v
        }
    }

    return merged
}
