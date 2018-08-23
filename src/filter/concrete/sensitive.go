package filter_concrete

import "github.com/importcjj/sensitive"

type Sensitive struct {
	filter *sensitive.Filter
}

func (dc *Sensitive) Handle(body string, isCache bool, isStatic bool) string {
	if !isCache {
		body = dc.filter.Filter(body)
	}

	return body
}

func (dc *Sensitive) IsRequest() bool {
	return false
}

func (dc *Sensitive) Preload() {
	dc.filter = sensitive.New()
}
