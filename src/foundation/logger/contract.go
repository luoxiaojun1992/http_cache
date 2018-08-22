package logger

type loggerProto interface {
	Preload()
	Handle(err error)
}
