package extension

type extensionProto interface {
	StartUp()
	ShutDown()
	IsEnabled() int
}
