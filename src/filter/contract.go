package filter

import "net/http"

type filterProto interface {
	IsRequest() bool //Whether Request Filter
}

type requestFilterProto interface {
	Next(h http.Handler)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type responseFilterProto interface {
	Handle(body string) string
}
