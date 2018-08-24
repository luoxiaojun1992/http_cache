package server

import (
	"errors"
	"github.com/luoxiaojun1992/http_cache/src/cache"
	"github.com/luoxiaojun1992/http_cache/src/filter"
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
	"github.com/luoxiaojun1992/http_cache/src/foundation/logger"
	"github.com/luoxiaojun1992/http_cache/src/foundation/util"
	"github.com/luoxiaojun1992/http_cache/src/router"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
	"os/signal"
	"syscall"
	"os"
	"sync"
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
	routerConfig, err := router.FetchConfig(r.Host, uri)
	if err != nil {
		w.Write([]byte{})
		return
	}

	//Read Cache
	cacheKey := ""
	if r.Method == "GET" && routerConfig["cache"] == cache.ENABLED {
		cacheKey = routerConfig["host"] + uri
		multiCache := cache.MGetCache([]string{"header:" + cacheKey, "body:" + cacheKey})
		headerStr := multiCache[0]
		bodyStr := multiCache[1]
		if len(headerStr) > 0 {
			headers := util.DeSerialize(headerStr)
			for key, value := range headers {
				w.Header().Add(key, value)
			}
			if len(bodyStr) > 0 {
				w.Write([]byte(filter.OnResponse(bodyStr, true, false)))
				return
			}
		}
		logger.Do(errors.New("Cache Miss"))
	}

	//Proxy Request
	proxyR, err := http.NewRequest(r.Method, routerConfig["host"]+uri, r.Body)
	if err != nil {
		logger.Do(err)
		w.Write([]byte{})
		return
	}
	proxyR.Header = r.Header
	urlObj, err := url.Parse(routerConfig["host"])
	if err == nil {
		proxyR.Header.Add("Host", urlObj.Host)
	}
	for _, cookie := range r.Cookies() {
		proxyR.AddCookie(cookie)
	}
	client := &http.Client{}
	resp, err := client.Do(proxyR)
	if err != nil {
		logger.Do(err)
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
	if r.Method == "GET" && routerConfig["cache"] == cache.ENABLED {
		ttl, err := time.ParseDuration(routerConfig["ttl"])
		if err == nil {
			cache.SetCache("header:"+cacheKey, util.Serialize(resp.Header), ttl)
		}
	}

	//Transfer Body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Do(err)
		w.Write([]byte{})
		return
	}
	bodyStr := string(body)
	w.Write([]byte(filter.OnResponse(bodyStr, false, false)))

	//Update Body Cache
	if r.Method == "GET" && routerConfig["cache"] == cache.ENABLED {
		ttl, err := time.ParseDuration(routerConfig["ttl"])
		if err == nil {
			cache.SetCache("body:"+cacheKey, filter.OnResponse(bodyStr, false, true), ttl)
		}
	}
}

func StartHttp() {
	var wg sync.WaitGroup

	s := &http.Server{
		Addr:           Env("HTTP_HOST", "0.0.0.0") + ":" + Env("HTTP_PORT", "8888"),
		Handler:        filter.OnRequest(&myHandler{}),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	wg.Add(1)
	go func(){
		log.Println(s.ListenAndServe())
		wg.Done()
	}()

	//Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	//Graceful Shutdown
	s.Shutdown(nil)

	wg.Wait()
}
