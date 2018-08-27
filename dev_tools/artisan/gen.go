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

import "net/http"

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

func ({shortName} *{name}) IsEnabled() int {
	return 0
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

func ({shortName} *{name}) IsEnabled() int {
	return 0
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
		log.Fatal("Name cannot be empty.")
	}

	if module < 0 {
		log.Fatal("Module cannot be less than zero.")
	}

	switch module {
	case LOGGER:
		genLogger(name)
	case RequestFilter:
		genRequestFilter(name)
	case ResponseFilter:
		genResponseFilter(name)
	default:
		log.Fatal("Unsupported module.")
	}
}

func parseName(name string) (lowerName, shortName string) {
	//todo handle camel name

	lowerName = strings.ToLower(name)
	shortName = lowerName[:1]
	return
}

func genLogger(name string) {
	genCode(name, "../../src/foundation/logger/concrete/", LoggerTpl)
}

func genRequestFilter(name string) {
	genCode(name, "../../src/filter/concrete/", RequestFilterTpl)
}

func genResponseFilter(name string) {
	genCode(name, "../../src/filter/concrete/", ResponseFilterTpl)
}

func genCode(name, dir, tpl string) {
	lowerName, shortName := parseName(name)
	filePath := dir + lowerName + ".go"
	tpl = strings.Replace(tpl, "{name}", name, -1)
	tpl = strings.Replace(tpl, "{shortName}", shortName, -1)
	err := writeFile(filePath, tpl)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Generated successfully.")
}

func writeFile(filePath string, fileContent string) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err == nil {
		defer file.Close()

		_, err := file.Write([]byte(fileContent))

		return err
	}

	return err
}
