package logger

type loggerProto interface {
	Preload()
	Error(err error)
	IsEnabled() int
	Warning(content string)
	Info(content string)
	Fatal(err error)
	Debug(content string)
	Trace(content string)
}
