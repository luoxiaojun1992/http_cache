package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"fmt"
)

type myHandler struct{}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Debug
	fmt.Println(r.URL.Path + "?" + r.URL.RawQuery + "#" + r.URL.Fragment)

	//Compose URI
	uri := r.URL.Path
	if len(r.URL.RawQuery) > 0 {
		uri += ("?" + r.URL.RawQuery)
	}
	if len(r.URL.Fragment) > 0 {
		uri += ("#" + r.URL.Fragment)
	}

	//Proxy Request
	resp, err := http.Get("https://www.baidu.com/" + uri)
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
	s := &http.Server{
		Addr:           ":8888",
		Handler:        &myHandler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
