package filter_concrete

import (
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
	"net/http"
)

type Header struct {
	next http.Handler
}

func (hd *Header) Next(h http.Handler) {
	hd.next = h
}

func (hd *Header) IsRequest() bool {
	return true
}

func (hd *Header) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hd.next.ServeHTTP(w, r)
}

func (hd *Header) IsEnabled() int {
	return EnvInt("HEADER_FILTER_SWITCH", 0)
}
