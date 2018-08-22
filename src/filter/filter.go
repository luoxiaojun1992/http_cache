package filter

import (
	. "github.com/luoxiaojun1992/http_cache/src/filter/concrete"
)

var filters []filter_proto

func Do(body string) string {
	for _, filter_concrete := range filters {
		body = filter_concrete.Handle(body)
	}

	return body
}

func InitFilter() {
	filters = []filter_proto{&DynamicContent{}}
}
