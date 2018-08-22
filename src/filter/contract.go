package filter

type filterProto interface {
	Handle(body string) string
}
