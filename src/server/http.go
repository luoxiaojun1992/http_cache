package server

import (
	"errors"
	"fmt"
	"github.com/luoxiaojun1992/http_cache/src/cache"
	. "github.com/luoxiaojun1992/http_cache/src/environment"
	"github.com/luoxiaojun1992/http_cache/src/filter"
	"github.com/luoxiaojun1992/http_cache/src/foundation/util"
	"github.com/luoxiaojun1992/http_cache/src/logger"
	"github.com/luoxiaojun1992/http_cache/src/router"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

//HTTP Handler
type myHandler struct{}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Compose URI
	uri := r.URL.RequestURI()
	if len(r.URL.Fragment) > 0 {
		uri += ("#" + r.URL.Fragment)
	}

	//Fetch Router Config
	router_config, err := router.FetchConfig(r.Host, uri)
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
			headers := util.DeSerialize(header_str)
			for key, value := range headers {
				w.Header().Add(key, value)
			}
			if len(body_str) > 0 {
				w.Write([]byte(filter.Do(body_str)))
				return
			}
		}
		logger.Do(errors.New("Cache Miss"))
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
	url_obj, err := url.Parse(router_config["host"])
	if err == nil {
		proxy_r.Header.Add("Host", url_obj.Host)
	}
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
		ttl, err := time.ParseDuration(router_config["ttl"])
		if err == nil {
			cache.SetCache("header:"+cache_key, util.Serialize(resp.Header), ttl)
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

func StartHttp() {
	s := &http.Server{
		Addr:           Env("HTTP_HOST", "0.0.0.0") + ":" + Env("HTTP_PORT", "8888"),
		Handler:        &myHandler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
