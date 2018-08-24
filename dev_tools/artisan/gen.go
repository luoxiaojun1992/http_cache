package main

import (
	"flag"
	"fmt"
	"github.com/ian-kent/go-log/log"
	"os"
	"strings"
)

const (
	LOGGER = iota
	RequestFilter
	ResponseFilter
)

const LoggerTpl = `package logger_concrete

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

const RequestFilterTpl = `package filter_concrete

type {name} struct {
	next http.Handler
}

func ({shortName} *{name}) Next(h http.Handler) {
	{shortName}.next = h
}

func ({shortName} *{name}) IsRequest() bool {
	return true
}

func ({shortName} *{name}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	{shortName}.next.ServeHTTP(w, r)
}
`

const ResponseFilterTpl = `package filter_concrete

type {name} struct {
}

func ({shortName} *{name}) Handle(body string, isCache bool, isStatic bool) string {
	return body
}

func ({shortName} *{name}) IsRequest() bool {
	return false
}

func ({shortName} *{name}) Preload() {
}
`

var name string
var module int

func init() {
	flag.StringVar(&name, "name", "", "Struct Name")
	flag.IntVar(&module, "module", -1, "Module")
	flag.Usage = func() {
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}
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
		genRequestFilter(name)
	case ResponseFilter:
		genResponseFilter(name)
	}
}

func parseName(name string) (lowerName, shortName string) {
	lowerName = strings.ToLower(name)
	shortName = lowerName[:1]
	return
}

func genLogger(name string) {
	lowerName, shortName := parseName(name)
	filePath := "../../src/foundation/logger/concrete/" + lowerName + ".go"
	loggerTpl := strings.Replace(LoggerTpl, "{name}", name, -1)
	loggerTpl = strings.Replace(loggerTpl, "{shortName}", shortName, -1)
	writeFile(filePath, loggerTpl)
}

func genRequestFilter(name string) {
	lowerName, shortName := parseName(name)
	filePath := "../../src/filter/concrete/" + lowerName + ".go"
	loggerTpl := strings.Replace(RequestFilterTpl, "{name}", name, -1)
	loggerTpl = strings.Replace(loggerTpl, "{shortName}", shortName, -1)
	writeFile(filePath, loggerTpl)
}

func genResponseFilter(name string) {
	lowerName, shortName := parseName(name)
	filePath := "../../src/filter/concrete/" + lowerName + ".go"
	loggerTpl := strings.Replace(ResponseFilterTpl, "{name}", name, -1)
	loggerTpl = strings.Replace(loggerTpl, "{shortName}", shortName, -1)
	writeFile(filePath, loggerTpl)
}

func writeFile(filePath string, fileContent string) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err == nil {
		defer file.Close()

		file.Write([]byte(fileContent))
	}
}
