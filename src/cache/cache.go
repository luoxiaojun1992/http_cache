package cache

import (
	"github.com/go-redis/redis"
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	ENABLED = "1"
)

type myCache struct {
	localCache       *cache.Cache
	localCacheSwitch int
	redisClient      *redis.Client
	prefix           string
}

var cacheObj *myCache

func NewCache() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     Env("REDIS_HOST", "localhost") + ":" + Env("REDIS_PORT", "6379"),
		Password: Env("REDIS_PASSWORD", ""), // no password set
		DB:       EnvInt("REDIS_DB", 0),     // use default DB
		PoolSize: EnvInt("REDIS_POOL_SIZE", 200),
	})

	cacheObj = &myCache{
		localCache:       cache.New(1*time.Second, 10*time.Minute),
		localCacheSwitch: EnvInt("LOCAL_CACHE_SWITCH", 0),
		redisClient:      redisClient,
		prefix:           Env("CACHE_PREFIX", ""),
	}
}

func (c *myCache) close() {
	c.redisClient.Close()
}

func (c *myCache) setCache(key, value string, ttl time.Duration) {
	if c.localCacheSwitch == 0 {
		return
	}
	c.localCache.Add(c.prefix+key, value, ttl)
}

func (c *myCache) getCache(key string) string {
	if c.localCacheSwitch == 0 {
		return ""
	}

	if x, found := c.localCache.Get(c.prefix + key); found {
		return x.(string)
	}

	return ""
}

func (c *myCache) setRedis(key, value string, ttl time.Duration) {
	c.redisClient.Set(c.prefix+key, value, ttl)
}

func (c *myCache) getRedis(key string) string {
	val, err := c.redisClient.Get(c.prefix + key).Result()
	if err == nil {
		return val
	}

	return ""
}

func Close() {
	cacheObj.close()
}

func SetCache(key, value string, ttl time.Duration) {
	cacheObj.setRedis(key, value, ttl)
	SetLocalCache(key, value)
}

func SetLocalCache(key, value string) {
	cacheObj.setCache(key, value, 1*time.Second)
}

func MGetCache(keys []string) []string {
	var vals []string
	for _, key := range keys {
		localCache := cacheObj.getCache(key)
		if len(localCache) <= 0 {
			localCache = cacheObj.getRedis(key)
			if len(localCache) > 0 {
				SetLocalCache(key, localCache)
			}
		}
		vals = append(vals, localCache)
	}
	return vals
}
