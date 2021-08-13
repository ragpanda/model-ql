package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

type LoggerWrapper struct {
	logger *logrus.Logger
}

var logger *LoggerWrapper
var onceLogger sync.Once

func Logger() *LoggerWrapper {
	onceLogger.Do(func() {
		logger = &LoggerWrapper{
			logger: logrus.StandardLogger(),
		}
		logger.logger.Level = logrus.DebugLevel
		logger.logger.SetFormatter(&CustomerFormatter{})
		logger.logger.AddHook(NewContextHook(nil))
	})
	return logger
}

func Debug(ctx context.Context, fmtStr string, args ...interface{}) {
	Logger().logger.WithContext(ctx).Debugf(fmtStr, args...)
}

func Info(ctx context.Context, fmtStr string, args ...interface{}) {
	Logger().logger.WithContext(ctx).Infof(fmtStr, args...)
}

func Warn(ctx context.Context, fmtStr string, args ...interface{}) {
	Logger().logger.WithContext(ctx).Warnf(fmtStr, args...)
}

func Error(ctx context.Context, fmtStr string, args ...interface{}) {
	Logger().logger.WithContext(ctx).Errorf(fmtStr, args...)
}

func Fatal(ctx context.Context, fmtStr string, args ...interface{}) {
	Logger().logger.WithContext(ctx).Fatalf(fmtStr, args...)
}

type CustomerFormatter struct {
}

// logTypeToColor converts the Level to a color string.
func (self *CustomerFormatter) logTypeToColor(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "[0;37"
	case logrus.InfoLevel:
		return "[0;36"
	case logrus.WarnLevel:
		return "[0;33"
	case logrus.ErrorLevel:
		return "[0;31"
	case logrus.FatalLevel:
		return "[0;31"
	case logrus.PanicLevel:
		return "[0;31"
	}

	return "[0;37"
}

const (
	defaultLogTimeFormat = "2006-01-02 15:04:05.000"
)

// Formatter implements logrus.Formatter
func (self *CustomerFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	colorStr := self.logTypeToColor(entry.Level)
	fmt.Fprintf(b, "\033%sm", colorStr)

	fmt.Fprintf(b, "[%s] ", entry.Time.Format(defaultLogTimeFormat))

	if file, ok := entry.Data["file"]; ok {
		fmt.Fprintf(b, "[%s:%v] ", file, entry.Data["line"])
	}

	fmt.Fprintf(b, "[%s] %s", entry.Level.String(), entry.Message)

	b.WriteByte('\n')

	b.WriteString("\033[0m")

	return b.Bytes(), nil
}

type contextHook struct {
	SkipPkg []string
}

type InitContextParams struct {
}

func NewContextHook(skipPkg []string) *contextHook {
	skipPkg = append([]string{
		"github.com/sirupsen/logrus",
		"github.com/ragpanda/model-ql/util",
	})
	ch := &contextHook{
		SkipPkg: skipPkg,
	}

	return ch
}

// Fire implements logrus.Hook interface
// https://github.com/sirupsen/logrus/issues/63
func (hook *contextHook) Fire(entry *logrus.Entry) error {
	pc := make([]uintptr, 10)
	cnt := runtime.Callers(6, pc)

	entry.Data["file"] = ""
	entry.Data["line"] = 0
	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !hook.isSkippedPackageName(name) {
			file, line := fu.FileLine(pc[i] - 1)
			pathSp := strings.Split(file, "/")
			if len(pathSp) > 3 {
				file = path.Join(pathSp[len(pathSp)-3:]...)
			}
			entry.Data["file"] = file
			entry.Data["line"] = line
			break
		}
	}

	return nil
}
func (hook *contextHook) isSkippedPackageName(name string) bool {
	for _, pkgName := range hook.SkipPkg {
		if strings.Contains(name, pkgName) {
			return true
		}
	}
	return false

}

func show(data interface{}) string {
	dataValue := reflect.Indirect(reflect.ValueOf(data))
	kind := dataValue.Type().Kind()
	if kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Map {
		b, _ := json.Marshal(data)
		return string(b)

	} else {
		result, err := cast.ToStringE(data)
		if err != nil {
			return fmt.Sprintf("%+v", result)
		}
		return result
	}

}

// Levels implements logrus.Hook interface.
func (hook *contextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
