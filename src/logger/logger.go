package logger

var loggers []logger_proto

func InitLogger() {
	loggers = append(loggers, &sentry{})
	for _, logger_concrete := range loggers {
		logger_concrete.preload()
	}
}

func Do(err error) {
	for _, logger_concrete := range loggers {
		logger_concrete.handle(err)
	}
}
