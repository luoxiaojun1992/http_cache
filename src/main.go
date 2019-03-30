package main

import (
	"github.com/luoxiaojun1992/http_cache/src/cache"
	"github.com/luoxiaojun1992/http_cache/src/filter"
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
	"github.com/luoxiaojun1992/http_cache/src/foundation/extension"
	"github.com/luoxiaojun1992/http_cache/src/foundation/logger"
	"github.com/luoxiaojun1992/http_cache/src/redis"
	"github.com/luoxiaojun1992/http_cache/src/router"
	"github.com/luoxiaojun1992/http_cache/src/server"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func init() {
	InitEnv()
	logger.InitLogger()
	router.InitConfig()
	cache.InitCache()
	redis.InitRedis()
	filter.InitFilter()
	extension.StartUp()
}

func main() {
	//pprof
	if EnvInt("PPROF_SWITCH", 0) == 1 {
		go func() {
			log.Println(http.ListenAndServe(Env("PPROF_HOST", "localhost")+":"+Env("PPROF_PORT", "6060"), nil))
		}()
	}

	defer redis.Close()
	defer extension.ShutDown()

	server.StartHttp()
}
