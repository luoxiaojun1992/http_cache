package filter_concrete

import (
	"github.com/importcjj/sensitive"
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
)

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

func (s *Sensitive) IsEnabled() int {
	return EnvInt("SENSITIVE_FILTER_SWITCH", 0)
}
