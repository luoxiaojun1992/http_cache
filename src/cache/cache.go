package cache

import (
	"github.com/go-redis/redis"
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
	"github.com/patrickmn/go-cache"
	"time"
	"strings"
	"github.com/luoxiaojun1992/http_cache/src/foundation/logger"
)

const (
	ENABLED  = "1"
	DISABLED = "0"
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

func (mc *myCache) close() {
	mc.redisClient.Close()
}

func (mc *myCache) setCache(key, value string, ttl time.Duration) {
	if mc.localCacheSwitch == 0 {
		return
	}
	mc.localCache.Set(mc.prefix+key, value, ttl)
}

func (mc *myCache) getCache(key string) string {
	if mc.localCacheSwitch == 0 {
		return ""
	}

	if x, found := mc.localCache.Get(mc.prefix + key); found {
		return x.(string)
	}

	return ""
}

func (mc *myCache) incrementCache(key string, step int, ttl time.Duration) int {
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

func (mc *myCache) setRedis(key, value string, ttl time.Duration) {
	mc.redisClient.Set(mc.prefix+key, value, ttl)
}

func (mc *myCache) getRedis(key string) string {
	val, err := mc.redisClient.Get(mc.prefix + key).Result()
	if err == nil {
		return val
	} else {
		logger.Error(err)
	}

	return ""
}

func (mc *myCache) delRedis(key string) (int64, error) {
	if strings.Contains(key, "*") {
		var cursor uint64
		var deleted int64
		for {
			var keys []string
			var err error
			keys, cursor, err = mc.redisClient.Scan(cursor, key, 10).Result()
			if err == nil {
				n, err := mc.redisClient.Del(keys ...).Result()
				if err != nil {
					logger.Error(err)
					return deleted, err
				}

				deleted += n

				if cursor == 0 {
					break
				}
			} else {
				logger.Error(err)
				return deleted, err
			}
		}

		return deleted, nil
	} else {
		deleted, err := mc.redisClient.Del(key).Result()
		if err != nil {
			logger.Error(err)
		}

		return deleted, err
	}
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
	var values []string
	for _, key := range keys {
		localCache := cacheObj.getCache(key)
		if len(localCache) <= 0 {
			localCache = cacheObj.getRedis(key)
			if len(localCache) > 0 {
				SetLocalCache(key, localCache)
			}
		}
		values = append(values, localCache)
	}
	return values
}

func IncrementLocalCache(key string, step int, ttl time.Duration) int {
	return cacheObj.incrementCache(key, step, ttl)
}

func DelRedis(key string) (int64, error) {
	return cacheObj.delRedis(key)
}
