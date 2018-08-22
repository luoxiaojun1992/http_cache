package logger

type logger_proto interface {
	Preload()
	Handle(err error)
}
