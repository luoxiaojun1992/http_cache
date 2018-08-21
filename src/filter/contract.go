package filter

type filter_proto interface {
	handle(body string) string
}
