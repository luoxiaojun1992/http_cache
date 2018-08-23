package filter_concrete

import "github.com/importcjj/sensitive"

type Sensitive struct {
	filter *sensitive.Filter
}

func (s *Sensitive) Handle(body string, isCache bool, isStatic bool) string {
	if !isCache {
		body = s.filter.Filter(body)
	}

	return body
}

func (s *Sensitive) IsRequest() bool {
	return false
}

func (s *Sensitive) Preload() {
	s.filter = sensitive.New()
}
