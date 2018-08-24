package main

import (
	"flag"
	"github.com/ian-kent/go-log/log"
	"strings"
	"os"
)

const (
	LOGGER = iota
	RequestFilter
	ResponseFilter
)

const LoggerTpl  = `package logger_concrete

type {name} struct {
}

func ({shortName} *{name}) Preload() {
}

func ({shortName} *{name}) Error(err error) {
}

func ({shortName} *{name}) Warning(content string) {
}

func ({shortName} *{name}) Info(content string) {
}

func ({shortName} *{name}) Fatal(err error) {
}

func ({shortName} *{name}) Debug(content string) {
}

func ({shortName} *{name}) Trace(content string) {
}

func ({shortName} *{name}) IsEnabled() int {
	return 0
}
`

var name string
var module int

func init() {
	flag.StringVar(&name, "name", "", "Struct Name")
	flag.IntVar(&module, "module", -1, "Module")
}

func main() {
	flag.Parse()

	if len(name) <= 0 {
		log.Fatal("Name cannot be empty")
	}

	if module < 0 {
		log.Fatal("Module cannot be less than zero")
	}

	switch module {
		case LOGGER:
			genLogger(name)
		case RequestFilter:
		case ResponseFilter:
	}
}

func genLogger(name string) {
	lowerName := strings.ToLower(name)
	filePath := "../../src/foundation/logger/concrete/" + lowerName + ".go"
	shortName := lowerName[:1]
	loggerTpl := strings.Replace(LoggerTpl, "{name}", name, -1)
	loggerTpl = strings.Replace(loggerTpl, "{shortName}", shortName, -1)
	writeFile(filePath, loggerTpl)
}

func genRequestFilter(name string) {
	//todo
}

func genResponseFilter(name string) {
	//todo
}

func writeFile(filePath string, fileContent string) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err == nil {
		defer file.Close()

		file.Write([]byte(fileContent))
	}
}
