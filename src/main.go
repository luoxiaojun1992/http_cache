package main

import (
	"github.com/luoxiaojun1992/http_cache/src/cache"
	. "github.com/luoxiaojun1992/http_cache/src/environment"
	"github.com/luoxiaojun1992/http_cache/src/filter"
	"github.com/luoxiaojun1992/http_cache/src/logger"
	"github.com/luoxiaojun1992/http_cache/src/router"
	"github.com/luoxiaojun1992/http_cache/src/server"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	//todo cache clean & hotfix
	//todo monitor http handler

	//Init Env
	InitEnv()

	//pprof
	if EnvInt("PPROF_SWITCH", 0) == 1 {
		go func() {
			log.Println(http.ListenAndServe(Env("PPROF_HOST", "localhost")+":"+Env("PPROF_PORT", "6060"), nil))
		}()
	}

	//Init Cache
	cache.NewCache()
	defer cache.Close()

	//Init Router Config
	router.InitConfig()

	//Init Filters
	filter.InitFilter()

	//Init Logger
	logger.InitLogger()

	//Start Proxy Server
	server.StartHttp()
}
