package cache

import (
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
	"github.com/luoxiaojun1992/http_cache/src/foundation/logger"
	"github.com/patrickmn/go-cache"
	"time"
)

type myCache struct {
	localCache       *cache.Cache
	localCacheSwitch int
	prefix           string
}

var cacheObj *myCache

// InitCache ...
func InitCache() {
	cacheObj = &myCache{
		localCache:       cache.New(1*time.Second, 10*time.Minute),
		localCacheSwitch: EnvInt("LOCAL_CACHE_SWITCH", 0),
		prefix:           Env("CACHE_PREFIX", ""),
	}
}

func (mc *myCache) set(key, value string, ttl time.Duration) {
	if mc.localCacheSwitch == 0 {
		return
	}
	mc.localCache.Set(mc.prefix+key, value, ttl)
}

func (mc *myCache) get(key string) string {
	if mc.localCacheSwitch == 0 {
		return ""
	}

	if x, found := mc.localCache.Get(mc.prefix + key); found {
		return x.(string)
	}

	return ""
}

func (mc *myCache) incrementBy(key string, step int, ttl time.Duration) int {
	key = mc.prefix + key

	err := mc.localCache.Add(key, step, ttl)
	if err == nil {
		return step
	} else {
		logger.Error(err)
	}

	newValue, err := mc.localCache.IncrementInt(key, step)
	if err == nil {
		return newValue
	} else {
		logger.Error(err)
	}

	return 0
}

// Cache set ...
func Set(key, value string, ttl time.Duration) {
	cacheObj.set(key, value, ttl)
}

// Cache get ...
func Get(key string) string {
	return cacheObj.get(key)
}

// Cache IncrementBy ...
func IncrementLocalCache(key string, step int, ttl time.Duration) int {
	return cacheObj.incrementBy(key, step, ttl)
}
