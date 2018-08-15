package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	CACHE_ENABLED = "1"
)

//ThirdParty Json Searilizer
var json = jsoniter.ConfigCompatibleWithStandardLibrary

//Cache Prefix
var cache_prefix string

//Cache Storage
var redis_client *redis.Client

//Router Config
var router map[string](map[string](map[string]string))

//HTTP Handler
type myHandler struct{}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Compose URI
	uri := r.URL.RequestURI()
	if len(r.URL.Fragment) > 0 {
		uri += r.URL.Fragment
	}

	//Fetch Router Config
	request_host := r.Header.Get("x-real-host")
	router_config := make(map[string]string)
	v, ok := router[request_host][uri]
	if !ok {
		v, ok := router[request_host]["*"]
		if !ok {
			w.Write([]byte{})
			return
		} else {
			router_config = v
		}
	} else {
		router_config = v
	}

	//Read Cache
	cache_key := ""
	if r.Method == "GET" && router_config["cache"] == CACHE_ENABLED {
		cache_key = cache_prefix + router_config["host"] + uri
		//todo cache key prefix
		result_arr, err := redis_client.MGet("header:"+cache_key, "body:"+cache_key).Result()
		if err == nil {
			if header_str, ok := result_arr[0].(string); ok {
				if len(header_str) > 0 {
					val := make(map[string][]string)
					json.Unmarshal([]byte(header_str), &val)
					for key, values := range val {
						for _, value := range values {
							w.Header().Add(key, value)
						}
					}
					if body_str, ok := result_arr[1].(string); ok {
						if len(body_str) > 0 {
							w.Write([]byte(fillDynamicContent(body_str)))
							return
						}
					}
				}
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
	if r.Method == "GET" && router_config["cache"] == CACHE_ENABLED {
		header_str, err := json.Marshal(resp.Header)
		if err == nil {
			ttl, err := time.ParseDuration(router_config["ttl"])
			if err == nil {
				redis_client.Set("header:"+cache_key, string(header_str), ttl)
			}
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
	w.Write([]byte(fillDynamicContent(body_str)))

	//Update Body Cache
	if r.Method == "GET" && router_config["cache"] == CACHE_ENABLED {
		ttl, err := time.ParseDuration(router_config["ttl"])
		if err == nil {
			redis_client.Set("body:"+cache_key, body_str, ttl)
		}
	}
}

func main() {
	//todo cache clean & hotfix

	//pprof
	if envInt("PPROF_SWITCH", 0) == 1 {
		go func() {
			log.Println(http.ListenAndServe(env("PPROF_HOST", "localhost")+":"+env("PPROF_PORT", "6060"), nil))
		}()
	}

	//Init Env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//Init Cache Prefix
	cache_prefix = env("CACHE_PREFIX", "")

	//Init Cache
	redis_client = redis.NewClient(&redis.Options{
		Addr:     env("REDIS_HOST", "localhost") + ":" + env("REDIS_PORT", "6379"),
		Password: env("REDIS_PASSWORD", ""), // no password set
		DB:       envInt("REDIS_DB", 0),     // use default DB
		PoolSize: envInt("REDIS_POOL_SIZE", 200),
	})
	defer redis_client.Close()

	//Router Config
	router = make(map[string](map[string](map[string]string)))
	router_config, err := ioutil.ReadFile("../router_config.json")
	if err != nil {
		log.Fatal(err)
	}
	if len(router_config) > 0 {
		json.Unmarshal(router_config, &router)
	}

	//Start Proxy Server
	s := &http.Server{
		Addr:           env("HTTP_HOST", "0.0.0.0") + ":" + env("HTTP_PORT", "8888"),
		Handler:        &myHandler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func fillDynamicContent(body string) string {
	re := regexp.MustCompile(`\<dynamic\>.+\</dynamic\>`)
	dynamic_tags := re.FindAllString(body, -1)
	if len(dynamic_tags) <= 0 {
		return body
	}

	dynamic_contents := make(map[string]string)
	var wg sync.WaitGroup
	var mutex_lock sync.Mutex
	for i, dynamic_tag := range dynamic_tags {
		if i >= 10 {
			break
		}
		go func() {
			defer wg.Done()

			dynamic_url := strings.Replace(dynamic_tag, "<dynamic>", "", 1)
			dynamic_url = strings.Replace(dynamic_url, "</dynamic>", "", 1)
			if val, ok := dynamic_contents[dynamic_url]; ok {
				mutex_lock.Lock()
				body = strings.Replace(body, dynamic_tag, val, 1)
				mutex_lock.Unlock()
				return
			}

			resp, err := http.Get(dynamic_url)
			if err == nil {
				defer resp.Body.Close()

				dynamic_content, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					dynamic_content_str := string(dynamic_content)
					dynamic_contents[dynamic_url] = dynamic_content_str
					mutex_lock.Lock()
					body = strings.Replace(body, dynamic_tag, dynamic_content_str, 1)
					mutex_lock.Unlock()
				}
			}
		}()
		wg.Add(1)
	}
	wg.Wait()

	return body
}

func env(key, default_value string) string {
	val := os.Getenv(key)

	if len(val) > 0 {
		return val
	}

	return default_value
}

func envInt(key string, default_value int) int {
	val := env(key, "")
	if len(val) > 0 {
		i, err := strconv.Atoi(val)
		if err == nil {
			return i
		}
	}

	return default_value
}
