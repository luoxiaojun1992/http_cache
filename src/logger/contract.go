package logger

type logger_proto interface {
	preload()
	handle(err error)
}
