package filter

var filters []filter_proto

func Do(body string) string {
	for _, filter_concrete := range filters {
		body = filter_concrete.handle(body)
	}

	return body
}

func InitFilter() {
	filters = []filter_proto{&dynamic_content{}}
}
