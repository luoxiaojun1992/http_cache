package filter

type filter_proto interface {
	Handle(body string) string
}
