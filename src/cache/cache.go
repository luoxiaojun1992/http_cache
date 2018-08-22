package cache

import (
	"github.com/go-redis/redis"
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	CACHE_ENABLED = "1"
)

type mycache struct {
	local_cache        *cache.Cache
	local_cache_switch int
	redis_client       *redis.Client
	prefix             string
}

var cache_obj *mycache

func NewCache() {
	redis_client := redis.NewClient(&redis.Options{
		Addr:     Env("REDIS_HOST", "localhost") + ":" + Env("REDIS_PORT", "6379"),
		Password: Env("REDIS_PASSWORD", ""), // no password set
		DB:       EnvInt("REDIS_DB", 0),     // use default DB
		PoolSize: EnvInt("REDIS_POOL_SIZE", 200),
	})

	cache_obj = &mycache{
		local_cache:        cache.New(1*time.Second, 10*time.Minute),
		local_cache_switch: EnvInt("LOCAL_CACHE_SWITCH", 0),
		redis_client:       redis_client,
		prefix:             Env("CACHE_PREFIX", ""),
	}
}

func (c *mycache) close() {
	c.redis_client.Close()
}

func (c *mycache) setCache(key, value string, ttl time.Duration) {
	if c.local_cache_switch == 0 {
		return
	}
	c.local_cache.Add(c.prefix+key, value, ttl)
}

func (c *mycache) getCache(key string) string {
	if c.local_cache_switch == 0 {
		return ""
	}

	if x, found := c.local_cache.Get(c.prefix + key); found {
		return x.(string)
	}

	return ""
}

func (c *mycache) setRedis(key, value string, ttl time.Duration) {
	c.redis_client.Set(c.prefix+key, value, ttl)
}

func (c *mycache) getRedis(key string) string {
	val, err := c.redis_client.Get(c.prefix + key).Result()
	if err == nil {
		return val
	}

	return ""
}

func Close() {
	cache_obj.close()
}

func SetCache(key, value string, ttl time.Duration) {
	cache_obj.setRedis(key, value, ttl)
	SetLocalCache(key, value)
}

func SetLocalCache(key, value string) {
	cache_obj.setCache(key, value, 1*time.Second)
}

func MGetCache(keys []string) []string {
	vals := []string{}
	for _, key := range keys {
		local_cache := cache_obj.getCache(key)
		if len(local_cache) <= 0 {
			local_cache = cache_obj.getRedis(key)
			if len(local_cache) > 0 {
				SetLocalCache(key, local_cache)
			}
		}
		vals = append(vals, local_cache)
	}
	return vals
}
