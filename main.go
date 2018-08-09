package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var router map[string]string

type myHandler struct{}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Proxy Request
	uri := r.URL.RequestURI()
	if len(r.URL.Fragment) > 0 {
		uri += r.URL.Fragment
	}
	proxy_r, err := http.NewRequest(r.Method, router["www.baidu.com"]+uri, r.Body)
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

	//Transfer Body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//todo log
		fmt.Println(err)
		w.Write([]byte{})
		return
	}
	w.Write(body)
}

func main() {
	//Router Config
	router = make(map[string]string)
	router["www.baidu.com"] = "http://www.dodoca.com"

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
