package logger

import (
	. "github.com/luoxiaojun1992/http_cache/src/foundation/logger/concrete"
	"github.com/pkg/errors"
)

const (
	ERROR = iota
	WARNING
	INFO
	FATAL
	DEBUG
	TRACE
)

var loggers []loggerProto

func InitLogger() {
	loggers = []loggerProto{&Sentry{}, &File{}}
	for _, loggerConcrete := range loggers {
		loggerConcrete.Preload()
	}
}

func Log(err error, level int) {
	for _, loggerConcrete := range loggers {
		if loggerConcrete.IsEnabled() == 0 {
			continue
		}

		switch level {
		case ERROR:
			loggerConcrete.Error(err)
		case WARNING:
			loggerConcrete.Warning(err.Error())
		case FATAL:
			loggerConcrete.Fatal(err)
		case DEBUG:
			loggerConcrete.Debug(err.Error())
		case TRACE:
			loggerConcrete.Trace(err.Error())
		case INFO:
			fallthrough
		default:
			loggerConcrete.Info(err.Error())
		}
	}
}

func Error(err error) {
	Log(err, ERROR)
}

func Warning(content string) {
	Log(errors.New(content), WARNING)
}

func Info(content string) {
	Log(errors.New(content), INFO)
}

func Fatal(err error) {
	Log(err, FATAL)
}

func Debug(content string) {
	Log(errors.New(content), DEBUG)
}

func Trace(content string) {
	Log(errors.New(content), TRACE)
}
