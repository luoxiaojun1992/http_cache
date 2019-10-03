package router

import (
	"errors"
	"github.com/json-iterator/go"
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
	"github.com/luoxiaojun1992/http_cache/src/foundation/logger"
	"io/ioutil"
)

const (
	CACHE_DISABLED = iota
	CACHE_ENABLED
)

//ThirdParty Json Serializer
var json = jsoniter.ConfigCompatibleWithStandardLibrary

//Router Config
var router map[string](map[string](map[string]string))

func InitConfig() {
	router = make(map[string](map[string](map[string]string)))
	routerConfig, err := ioutil.ReadFile(Env("ROUTER_CONFIG_FILE_PATH", "../etc/router_config.json"))
	if err != nil {
		logger.Fatal(err)
	}
	if len(routerConfig) > 0 {
		json.Unmarshal(routerConfig, &router)
	}
}

func FetchConfig(requestHost string, uri string) (map[string]string, error) {
	routerConfig := make(map[string]string)
	v, ok := router[requestHost][uri]
	if !ok {
		v, ok := router[requestHost]["*"]
		if !ok {
			v, ok := router["*"][uri]
			if !ok {
				v, ok := router["*"]["*"]
				if !ok {
					return routerConfig, errors.New("router config not set")
				} else {
					routerConfig = v
				}
			} else {
				routerConfig = v
			}
		} else {
			routerConfig = v
		}
	} else {
		routerConfig = v
	}

	return routerConfig, nil
}
