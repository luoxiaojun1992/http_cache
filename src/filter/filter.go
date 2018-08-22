package filter

import (
	. "github.com/luoxiaojun1992/http_cache/src/filter/concrete"
	"net/http"
)

var requestFilters []requestFilterProto
var responseFilters []responseFilterProto

func OnRequest(h http.Handler) http.Handler {
	requestFiltersLen := len(requestFilters)
	for i, filterConcrete := range requestFilters {
		if i > 0 {
			requestFilters[i-1].Next(filterConcrete.(http.Handler))
		}

		if i == requestFiltersLen-1 {
			filterConcrete.Next(h)
		}
	}

	return requestFilters[0].(http.Handler)
}

func OnResponse(body string) string {
	for _, filterConcrete := range responseFilters {
		body = filterConcrete.Handle(body)
	}

	return body
}

func InitFilter() {
	allFilters := []filterProto{&DynamicContent{}, &FlowControl{}}

	for _, filter := range allFilters {
		if filter.IsRequest() {
			requestFilters = append(requestFilters, filter.(requestFilterProto))
		} else {
			responseFilters = append(responseFilters, filter.(responseFilterProto))
		}
	}
}
