package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//Router Config
var router map[string]string

//Cache Storage
var body_cache map[string][]byte
var header_cache map[string]map[string][]string

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
	if val, ok := header_cache[cache_key]; ok {
		for key, values := range val {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
	}
	if val, ok := body_cache[cache_key]; ok {
		w.Write(val)
		return
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
	header_cache[cache_key] = resp.Header

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
	body_cache[cache_key] = body
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
	//todo using redis
	body_cache = make(map[string][]byte)
	header_cache = make(map[string]map[string][]string)

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
