package logger_concrete

import (
	"github.com/getsentry/raven-go"
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
)

type Sentry struct {
}

func (s *Sentry) Preload() {
	raven.SetDSN(Env("SENTRY_DSN", ""))
}

func (s *Sentry) Handle(err error) {
	raven.CaptureError(err, nil)
}

func (s *Sentry) IsEnabled() int {
	return EnvInt("SENTRY_SWITCH", 0)
}
