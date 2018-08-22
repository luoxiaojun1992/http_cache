package logger

import (
	. "github.com/luoxiaojun1992/http_cache/src/foundation/logger/concrete"
)

var loggers []logger_proto

func InitLogger() {
	loggers = []logger_proto{&Sentry{}}
	for _, logger_concrete := range loggers {
		logger_concrete.Preload()
	}
}

func Do(err error) {
	for _, logger_concrete := range loggers {
		logger_concrete.Handle(err)
	}
}
