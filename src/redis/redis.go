package redis

import (
	"github.com/go-redis/redis"
	. "github.com/luoxiaojun1992/http_cache/src/foundation/environment"
	"github.com/luoxiaojun1992/http_cache/src/foundation/logger"
	"github.com/patrickmn/go-cache"
	"strings"
	"time"
)

type myRedis struct {
	localCache       *cache.Cache
	localCacheSwitch int
	redisClient      *redis.Client
	prefix           string
}

var redisObj *myRedis

// InitRedis ...
func InitRedis() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     Env("REDIS_HOST", "localhost") + ":" + Env("REDIS_PORT", "6379"),
		Password: Env("REDIS_PASSWORD", ""),
		DB:       EnvInt("REDIS_DB", 0),
		PoolSize: EnvInt("REDIS_POOL_SIZE", 200),
	})

	redisObj = &myRedis{
		redisClient:      redisClient,
		prefix:           Env("CACHE_PREFIX", ""),
	}
}

func (mc *myRedis) close() {
	mc.redisClient.Close()
}

func (mc *myRedis) set(key, value string, ttl time.Duration) {
	mc.redisClient.Set(mc.prefix+key, value, ttl)
}

func (mc *myRedis) setNx(key, value string, ttl time.Duration) bool {
	res, err := mc.redisClient.SetNX(key, value, ttl).Result()
	if err == nil {
		return res
	} else {
		logger.Error(err)
	}

	return false
}

func (mc *myRedis) get(key string) string {
	val, err := mc.redisClient.Get(mc.prefix + key).Result()
	if err == nil {
		return val
	} else {
		logger.Error(err)
	}

	return ""
}

func (mc *myRedis) del(key string) (int64, error) {
	if strings.Contains(key, "*") {
		var cursor uint64
		var deleted int64
		for {
			var keys []string
			var err error
			keys, cursor, err = mc.redisClient.Scan(cursor, key, 10).Result()
			if err == nil {
				n, err := mc.redisClient.Del(keys...).Result()
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
	}

	deleted, err := mc.redisClient.Del(key).Result()
	if err != nil {
		logger.Error(err)
	}

	return deleted, err
}

// Close ...
func Close() {
	redisObj.close()
}

// Redis Set ...
func Set(key, value string, ttl time.Duration) {
	redisObj.set(key, value, ttl)
}

// Redis SetNx ...
func SetNx(key, value string, ttl time.Duration) bool {
	return redisObj.setNx(key, value, ttl)
}

// Redis Get ...
func Get(key string) string {
	return redisObj.get(key)
}

// Redis Del ...
func Del(key string) (int64, error) {
	return redisObj.del(key)
}
