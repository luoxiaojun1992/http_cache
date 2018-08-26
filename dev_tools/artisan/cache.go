package main

import (
	"flag"
	"fmt"
	"github.com/ian-kent/go-log/log"
	"github.com/luoxiaojun1992/http_cache/src/cache"
)

var op string
var key string

func init() {
	flag.StringVar(&op, "op", "clean", "Operation")
	flag.StringVar(&key, "key", "", "Cache Key")
	flag.Usage = func() {
		fmt.Println("Usage:")
		flag.PrintDefaults()
	}

	//Init Cache
	cache.NewCache()
}

func main() {
	defer cache.Close()

	flag.Parse()

	if len(key) <= 0 {
		log.Fatal("Cache key cannot be empty.")
	}

	switch op {
	case "clean":
		clean(key)
	default:
		log.Fatal("Unsupported operation.")
	}
}

func clean(key string) {
	deleted, err := cache.DelRedis(key)

	log.Println(fmt.Sprintf("Cleaned %d keys successfully", deleted))

	if err != nil {
		log.Fatal(err)
	}
}
