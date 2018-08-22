package filter_concrete

import "net/http"

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
	//todo flow control

	fc.next.ServeHTTP(w, r)
}
