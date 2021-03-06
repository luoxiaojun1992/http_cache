package filter

import (
	. "github.com/luoxiaojun1992/http_cache/src/filter/concrete"
	"net/http"
)

var requestFilters []requestFilterProto
var responseFilters []responseFilterProto

func OnRequest(h http.Handler) http.Handler {
	requestFilters = append(requestFilters, h.(requestFilterProto))

	for i, filterConcrete := range requestFilters {
		if i > 0 {
			requestFilters[i-1].Next(filterConcrete.(http.Handler))
		}
	}

	return requestFilters[0].(http.Handler)
}

func OnResponse(body string, isCache bool, isStatic bool) string {
	for _, filterConcrete := range responseFilters {
		body = filterConcrete.Handle(body, isCache, isStatic)
	}

	return body
}

func InitFilter() {
	allFilters := []filterProto{&FlowControl{}, &Header{}, &DynamicContent{}, &Sensitive{}}

	for _, filter := range allFilters {
		if filter.IsEnabled() == 0 {
			continue
		}

		filter.Preload()

		if filter.IsRequest() {
			requestFilters = append(requestFilters, filter.(requestFilterProto))
		} else {
			responseFilters = append(responseFilters, filter.(responseFilterProto))
		}
	}
}
