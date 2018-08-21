package router

import (
	"errors"
	"github.com/json-iterator/go"
	. "github.com/luoxiaojun1992/http_cache/src/environment"
	"io/ioutil"
	"log"
)

//ThirdParty Json Searilizer
var json = jsoniter.ConfigCompatibleWithStandardLibrary

//Router Config
var router map[string](map[string](map[string]string))

func InitConfig() {
	router = make(map[string](map[string](map[string]string)))
	router_config, err := ioutil.ReadFile(Env("ROUTER_CONFIG_FILE_PATH", "../etc/router_config.json"))
	if err != nil {
		log.Fatal(err)
	}
	if len(router_config) > 0 {
		json.Unmarshal(router_config, &router)
	}
}

func FetchConfig(request_host string, uri string) (map[string]string, error) {
	router_config := make(map[string]string)
	v, ok := router[request_host][uri]
	if !ok {
		v, ok := router[request_host]["*"]
		if !ok {
			return router_config, errors.New("router config not set")
		} else {
			router_config = v
		}
	} else {
		router_config = v
	}

	return router_config, nil
}
