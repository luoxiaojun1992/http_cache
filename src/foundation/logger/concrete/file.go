package logger_concrete

import (
	"github.com/ian-kent/go-log/appenders"
	"github.com/ian-kent/go-log/layout"
	"github.com/ian-kent/go-log/levels"
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
	logPath := logDir + "error.log"

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

func (f *File) Handle(err error) {
	go func() { f.logger.Log(levels.ERROR, err) }()
}

func (f *File) IsEnabled() int {
	return EnvInt("FILE_LOG_SWITCH", 0)
}
