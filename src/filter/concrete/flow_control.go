package filter_concrete

import (
	"github.com/luoxiaojun1992/http_cache/src/cache"
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
	"net/http"
	"time"
)

type FlowControl struct {
	next http.Handler
}

func (fc *FlowControl) Next(h http.Handler) {
	fc.next = h
}

func (fc *FlowControl) IsRequest() bool {
	return true
}

func (fc *FlowControl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cache.IncrementLocalCache("http_request_count", 1, 1*time.Second) >
		EnvInt("REQUEST_LIMIT_COUNT", 30000) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte{})
		return
	}

	fc.next.ServeHTTP(w, r)
}
