package logger_concrete

import (
	"github.com/getsentry/raven-go"
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
	"os"
)

type Sentry struct {
}

func (s *Sentry) Preload() {
	raven.SetDSN(Env("SENTRY_DSN", ""))
}

func (s *Sentry) Error(err error) {
	raven.CaptureError(err, nil)
}

func (s *Sentry) Warning(content string) {
	packet := raven.NewPacket(content)
	packet.Level = raven.WARNING
	raven.Capture(packet, nil)
}

func (s *Sentry) Info(content string) {
	packet := raven.NewPacket(content)
	packet.Level = raven.INFO
	raven.Capture(packet, nil)
}

func (s *Sentry) Fatal(err error) {
	packet := raven.NewPacket(err.Error())
	packet.Level = raven.FATAL
	_, ch := raven.Capture(packet, nil)
	<-ch
	os.Exit(2)
}

func (s *Sentry) Debug(content string) {
	packet := raven.NewPacket(content)
	packet.Level = raven.DEBUG
	raven.Capture(packet, nil)
}

func (s *Sentry) Trace(content string) {
	s.Info(content)
}

func (s *Sentry) IsEnabled() int {
	return EnvInt("SENTRY_SWITCH", 0)
}
