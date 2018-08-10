package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Router Config
var router map[string]string

//HTTP Handler
type myHandler struct{}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Compose URI
	uri := r.URL.RequestURI()
	if len(r.URL.Fragment) > 0 {
		uri += r.URL.Fragment
	}

	//Read Cache
	redis_client := redis.NewClient(&redis.Options{
		Addr:     env("REDIS_HOST", "localhost") + ":" + env("REDIS_PORT", "6379"),
		Password: env("REDIS_PASSWORD", ""), // no password set
		DB:       envInt("REDIS_DB", 0),     // use default DB
	})
	defer redis_client.Close()

	cache_key := cacheKey(r.Method, router[r.Host]+uri, r.Body)
	res, err := redis_client.Exists("header:" + cache_key).Result()
	if err == nil {
		if res > 0 {
			res, err := redis_client.Get("header:" + cache_key).Result()
			if err == nil && len(res) > 0 {
				val := make(map[string][]string)
				json.Unmarshal([]byte(res), &val)
				for key, values := range val {
					for _, value := range values {
						w.Header().Add(key, value)
					}
				}
				res, err := redis_client.Exists("body:" + cache_key).Result()
				if err == nil {
					if res > 0 {
						res, err := redis_client.Get("body:" + cache_key).Result()
						if err == nil {
							w.Write([]byte(fillDynamicContent(res)))
							return
						}
					}
				}
			}
		}
	}

	//Proxy Request
	proxy_r, err := http.NewRequest(r.Method, router[r.Host]+uri, r.Body)
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
	defer resp.Body.Close()
	if err != nil {
		//todo log
		fmt.Println(err)
		w.Write([]byte{})
		return
	}

	//Transfer Headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	//Update Header Cache
	header_str, err := json.Marshal(resp.Header)
	if err == nil {
		redis_client.Set("header:"+cache_key, string(header_str), 0)
	}

	//Transfer Body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//todo log
		fmt.Println(err)
		w.Write([]byte{})
		return
	}
	w.Write([]byte(fillDynamicContent(string(body))))

	//Update Body Cache
	redis_client.Set("body:"+cache_key, string(body), 0)
}

func main() {
	//todo goroutine config

	//Init Env
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	//Router Config
	router = make(map[string]string)
	//todo more complex
	router_config, err := ioutil.ReadFile("./router_config.json")
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

func cacheKey(method string, url string, body io.ReadCloser) string {
	body_str := ""
	if body != nil {
		body_byte, err := ioutil.ReadAll(body)
		if err != nil {
			//todo log
			fmt.Println(err)
			return ""
		}
		body_str = string(body_byte)
	}

	hmd5 := md5.New()
	io.WriteString(hmd5, method+":"+url+":"+body_str)
	return string(hmd5.Sum(nil))
}

func fillDynamicContent(body string) string {
	dynamic_contents := make(map[string]string)
	re := regexp.MustCompile(`\<dynamic\>.+\</dynamic\>`)
	dynamic_tags := re.FindAllString(body, -1)
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
			defer resp.Body.Close()
			if err == nil {
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
