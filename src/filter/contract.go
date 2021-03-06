package filter

import "net/http"

type filterProto interface {
	IsRequest() bool //Whether Request Filter
	IsEnabled() int
	Preload()
}

type requestFilterProto interface {
	Next(h http.Handler)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type responseFilterProto interface {
	Handle(body string, isCache bool, isStatic bool) string
}
