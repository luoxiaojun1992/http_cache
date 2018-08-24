package logger_concrete

import (
	"github.com/ian-kent/go-log/appenders"
	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/log"
	"github.com/ian-kent/go-log/logger"
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
	stdLog "log"
	"os"
)

type File struct {
	logger logger.Logger
}

func (f *File) Preload() {
	logDir := Env("LOG_DIR", "../logs/")
	logPath := logDir + "app.log"

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		stdLog.Fatal(err)
	}

	f.logger = log.Logger()
	r := appenders.RollingFile(logPath, true)
	r.MaxBackupIndex = EnvInt("MAX_LOG_FILE", 10)
	r.SetLayout(layout.Pattern("[%p] %d %m"))
	f.logger.SetAppender(r)
}

func (f *File) Error(err error) {
	go func() { f.logger.Error(err) }()
}

func (f *File) Warning(content string) {
	go func() { f.logger.Warn(content) }()
}

func (f *File) Info(content string) {
	go func() { f.logger.Info(content) }()
}

func (f *File) Fatal(err error) {
	go func() { f.logger.Fatal(err) }()
}

func (f *File) Debug(content string) {
	go func() { f.logger.Debug(content) }()
}

func (f *File) Trace(content string) {
	go func() { f.logger.Trace(content) }()
}

func (f *File) IsEnabled() int {
	return EnvInt("FILE_LOG_SWITCH", 0)
}
