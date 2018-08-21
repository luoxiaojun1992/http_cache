package logger

import (
	"github.com/getsentry/raven-go"
	. "github.com/luoxiaojun1992/http_cache/src/environment"
)

type sentry struct {
}

func (s *sentry) preload() {
	raven.SetDSN(Env("SENTRY_DSN", ""))
}

func (s *sentry) handle(err error) {
	raven.CaptureError(err, nil)
}
