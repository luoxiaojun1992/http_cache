package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//Router Config
var router map[string]string

//Cache Storage
var redis_client *redis.Client

//HTTP Handler
type myHandler struct{}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Compose URI
	uri := r.URL.RequestURI()
	if len(r.URL.Fragment) > 0 {
		uri += r.URL.Fragment
	}

	//Read Cache
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
				res, err := redis_client.Exists(cache_key).Result()
				if err == nil {
					if res > 0 {
						res, err := redis_client.Get(cache_key).Result()
						if err == nil {
							w.Write([]byte(res))
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
	w.Write(body)

	//Update Body Cache
	redis_client.Set(cache_key, string(body), 0)
}

func main() {
	//todo goroutine config

	//Router Config
	router = make(map[string]string)
	//todo more complex
	router_config, err := ioutil.ReadFile("./router_config.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(router_config, &router)

	//Init Cache
	//todo config
	redis_client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//Start Proxy Server
	s := &http.Server{
		Addr:           ":8888",
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
