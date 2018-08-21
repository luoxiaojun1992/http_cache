package main

import (
	"fmt"
	"github.com/luoxiaojun1992/http_cache/src/cache"
	. "github.com/luoxiaojun1992/http_cache/src/environment"
	"github.com/luoxiaojun1992/http_cache/src/filter"
	"github.com/luoxiaojun1992/http_cache/src/router"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"
)

//HTTP Handler
type myHandler struct{}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Compose URI
	uri := r.URL.RequestURI()
	if len(r.URL.Fragment) > 0 {
		uri += r.URL.Fragment
	}

	//Fetch Router Config
	request_host := r.Header.Get("x-request-host")
	router_config, err := router.FetchConfig(request_host, uri)
	if err != nil {
		w.Write([]byte{})
		return
	}

	//Read Cache
	cache_key := ""
	if r.Method == "GET" && router_config["cache"] == cache.CACHE_ENABLED {
		cache_key = router_config["host"] + uri
		multi_cache := cache.MGetCache([]string{"header:" + cache_key, "body:" + cache_key})
		header_str := multi_cache[0]
		body_str := multi_cache[1]

		if len(header_str) > 0 {
			headers := strings.Split(header_str, "\r\n\r\n")
			for _, header := range headers {
				header_pair := strings.Split(header, "\r\n")
				w.Header().Add(header_pair[0], header_pair[1])
			}
			if len(body_str) > 0 {
				w.Write([]byte(filter.Do(body_str)))
				return
			}
		}
		fmt.Println("Cache Miss")
	}

	//Proxy Request
	proxy_r, err := http.NewRequest(r.Method, router_config["host"]+uri, r.Body)
	if err != nil {
		//todo log
		fmt.Println(err)
		w.Write([]byte{})
		return
	}
	proxy_r.Header = r.Header
	for _, cookie := range r.Cookies() {
		proxy_r.AddCookie(cookie)
	}
	client := &http.Client{}
	resp, err := client.Do(proxy_r)
	if err != nil {
		//todo log
		fmt.Println(err)
		w.Write([]byte{})
		return
	}
	defer resp.Body.Close()

	//Transfer Headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	//Update Header Cache
	if r.Method == "GET" && router_config["cache"] == cache.CACHE_ENABLED {
		headers := []string{}
		for key, values := range resp.Header {
			for _, value := range values {
				headers = append(headers, key+"\r\n"+value)
			}
		}
		ttl, err := time.ParseDuration(router_config["ttl"])
		if err == nil {
			header_str := strings.Join(headers, "\r\n\r\n")
			cache.SetCache("header:"+cache_key, header_str, ttl)
		}
	}

	//Transfer Body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//todo log
		fmt.Println(err)
		w.Write([]byte{})
		return
	}
	body_str := string(body)
	w.Write([]byte(filter.Do(body_str)))

	//Update Body Cache
	if r.Method == "GET" && router_config["cache"] == cache.CACHE_ENABLED {
		ttl, err := time.ParseDuration(router_config["ttl"])
		if err == nil {
			cache.SetCache("body:"+cache_key, body_str, ttl)
		}
	}
}

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
	router.InitConfig(Env("ROUTER_CONFIG_FILE_PATH", "../router_config.json"))

	//Init Filters
	filter.InitFilter()

	//Start Proxy Server
	s := &http.Server{
		Addr:           Env("HTTP_HOST", "0.0.0.0") + ":" + Env("HTTP_PORT", "8888"),
		Handler:        &myHandler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
