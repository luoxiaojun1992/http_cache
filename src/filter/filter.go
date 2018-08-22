package filter

import (
	. "github.com/luoxiaojun1992/http_cache/src/filter/concrete"
)

var filters []filterProto

func Do(body string) string {
	for _, filterConcrete := range filters {
		body = filterConcrete.Handle(body)
	}

	return body
}

func InitFilter() {
	filters = []filterProto{&DynamicContent{}}
}
