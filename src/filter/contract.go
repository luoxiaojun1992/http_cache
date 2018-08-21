package filter

type Filter interface {
	Handle(body string) string
}
