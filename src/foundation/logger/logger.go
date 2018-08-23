package logger

import (
	. "github.com/luoxiaojun1992/http_cache/src/foundation/logger/concrete"
)

var loggers []loggerProto

func InitLogger() {
	loggers = []loggerProto{&Sentry{}, &File{}}
	for _, loggerConcrete := range loggers {
		loggerConcrete.Preload()
	}
}

func Do(err error) {
	for _, loggerConcrete := range loggers {
		if loggerConcrete.IsEnabled() != 0 {
			loggerConcrete.Handle(err)
		}
	}
}
