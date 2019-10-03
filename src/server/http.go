package server

import (
	"github.com/luoxiaojun1992/http_cache/src/cache"
	"github.com/luoxiaojun1992/http_cache/src/filter"
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
	"github.com/luoxiaojun1992/http_cache/src/foundation/logger"
	"github.com/luoxiaojun1992/http_cache/src/foundation/util"
	"github.com/luoxiaojun1992/http_cache/src/redis"
	"github.com/luoxiaojun1992/http_cache/src/router"
	"io/ioutil"
	stdLog "log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

//HTTP Handler
type myHandler struct{
	next http.Handler
}

func (h *myHandler) parseUri(r *http.Request) string {
	uri := r.URL.RequestURI()
	if len(r.URL.Fragment) > 0 {
		queryString := "#" + r.URL.Fragment
		uri += queryString
	}

	return uri
}

func (h *myHandler) updateHeaderCache(cacheKey string, headers map[string][]string, ttlConfig string) {
	ttl, err := time.ParseDuration(ttlConfig)
	if err == nil {
		serializedHeaders := util.Serialize(headers)
		if ttl.Seconds() > 0 {
			cache.Set("header:"+cacheKey, serializedHeaders, ttl)
		} else {
			cache.Set("header:"+cacheKey, serializedHeaders, 1*time.Second)
		}
		redis.Set("header:"+cacheKey+":ttl", ttl.String(), ttl)
		redis.Set("header:"+cacheKey, serializedHeaders, 0)
	} else {
		logger.Error(err)
	}
}

func (h *myHandler) updateBodyCache(cacheKey string, bodyStr string, ttlConfig string) {
	ttl, err := time.ParseDuration(ttlConfig)
	if err == nil {
		filteredBody := filter.OnResponse(bodyStr, false, true)
		if ttl.Seconds() > 0 {
			cache.Set("body:"+cacheKey, filteredBody, ttl)
		} else {
			cache.Set("body:"+cacheKey, filteredBody, 5*time.Second)
		}
		redis.Set("body:"+cacheKey+":ttl", ttl.String(), ttl)
		redis.Set("body:"+cacheKey, filteredBody, 0)
	} else {
		logger.Error(err)
	}
}

func (h *myHandler) Next(nextH http.Handler) {
	h.next = nextH
}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Compose URI
	uri := h.parseUri(r)

	//Fetch Router Config
	routerConfig, err := router.FetchConfig(r.Host, uri)
	if err != nil {
		w.Write([]byte{})
		return
	}

	//Read Cache
	cacheKey := ""
	if r.Method == "GET" && routerConfig["cache"] == strconv.Itoa(router.CACHE_ENABLED) {
		cacheKey = routerConfig["host"] + uri
		headerStr := cache.Get("header:" + cacheKey)
		if len(headerStr) <= 0 {
			if len(redis.Get("header:"+cacheKey+":ttl")) <= 0 {
				if redis.SetNx("header:"+cacheKey+":update:lock", "1", 5*time.Second) {
					defer redis.Del("header:"+cacheKey+":update:lock")
					headerStr = ""
				} else {
					headerStr = redis.Get("header:" + cacheKey)
				}
			} else {
				headerStr = redis.Get("header:" + cacheKey)
			}

			if len(headerStr) > 0 {
				ttl, err := time.ParseDuration(routerConfig["ttl"])
				if err == nil {
					if ttl.Seconds() > 0 {
						cache.Set("header:"+cacheKey, headerStr, ttl)
					} else {
						cache.Set("header:"+cacheKey, headerStr, 5*time.Second)
					}
				}
			}
		}

		bodyStr := cache.Get("body:" + cacheKey)
		if len(bodyStr) <= 0 {
			if len(redis.Get("body:"+cacheKey+":ttl")) <= 0 {
				if redis.SetNx("body:"+cacheKey+":update:lock", "1", 5*time.Second) {
					defer redis.Del("body:"+cacheKey+":update:lock")
					bodyStr = ""
				} else {
					bodyStr = redis.Get("body:" + cacheKey)
				}
			} else {
				bodyStr = redis.Get("body:" + cacheKey)
			}

			if len(bodyStr) > 0 {
				ttl, err := time.ParseDuration(routerConfig["ttl"])
				if err == nil {
					if ttl.Seconds() > 0 {
						cache.Set("body:"+cacheKey, bodyStr, ttl)
					} else {
						cache.Set("body:"+cacheKey, bodyStr, 5*time.Second)
					}
				}
			}
		}

		if len(headerStr) > 0 && len(bodyStr) > 0 {
			headers := util.DeSerialize(headerStr)
			for key, value := range headers {
				w.Header().Add(key, value)
			}

			w.Header().Add("X-Proxy-Cache", "hit")
			w.Write([]byte(filter.OnResponse(bodyStr, true, false)))
			return
		}
	}

	w.Header().Add("X-Proxy-Cache", "miss")

	//Proxy Request
	proxyR, err := http.NewRequest(r.Method, routerConfig["host"]+uri, r.Body)
	if err != nil {
		logger.Error(err)
		w.Write([]byte{})
		return
	}
	proxyR.Header = r.Header
	if routerConfig["preserve_host"] != "1" {
		urlObj, err := url.Parse(routerConfig["host"])
		if err == nil {
			proxyR.Header.Add("Host", urlObj.Host)
		} else {
			logger.Error(err)
		}
	}
	for _, cookie := range r.Cookies() {
		proxyR.AddCookie(cookie)
	}
	timeout := time.Second * 5 //Default 5s
	if configTimeout, errTimeout := time.ParseDuration(routerConfig["timeout"]); errTimeout == nil {
		timeout = configTimeout
	} else {
		logger.Error(errTimeout)
	}
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(proxyR)
	if err != nil {
		logger.Error(err)
		w.Write([]byte{})
		return
	}
	defer resp.Body.Close()

	//Cache Control Header
	cacheControl := ""

	//Transfer Headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)

			if strings.ToLower(key) == "cache-control" {
				cacheControl = strings.ToLower(value)
			}
		}
	}

	//Transfer HTTP Status Code
	w.WriteHeader(resp.StatusCode)

	//Transfer Body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		w.Write([]byte{})
		return
	}
	bodyStr := string(body)
	w.Write([]byte(filter.OnResponse(bodyStr, false, false)))

	//Determine if cache by http status code and cache control header
	if resp.StatusCode == http.StatusOK && util.IfCache(cacheControl) {
		if r.Method == "GET" && routerConfig["cache"] == strconv.Itoa(router.CACHE_ENABLED) {
			h.updateHeaderCache(cacheKey, resp.Header, routerConfig["ttl"])
		}

		//Update Body Cache
		if r.Method == "GET" && routerConfig["cache"] == strconv.Itoa(router.CACHE_ENABLED) {
			h.updateBodyCache(cacheKey, bodyStr, routerConfig["ttl"])
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
	go func() {
		stdLog.Println(s.ListenAndServe())
		wg.Done()
	}()

	stdLog.Println("Server started.")

	//Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	//Graceful Shutdown
	s.Shutdown(nil)

	wg.Wait()
}
