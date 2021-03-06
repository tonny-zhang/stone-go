package store

import (
	"stone/logger"
	"sync"
	"time"
)

var loggerCache = logger.GetPrefixLogger("storeCache")
var mu sync.Mutex

type cacheValue struct {
	time time.Time
	data map[string]interface{}
}

var cache map[string]cacheValue

func initCache() {
	if cache == nil {
		cache = make(map[string]cacheValue)
	}
}

// SetCache set cache value
func SetCache(key string, data map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()
	initCache()

	timeDelay := time.Now()
	delay, _ := time.ParseDuration("1h")
	timeDelay.Add(delay)
	cache[key] = cacheValue{
		time: timeDelay,
		data: data,
	}
}

// GetCache get value from cache
func GetCache(key string) map[string]interface{} {
	mu.Lock()
	defer mu.Unlock()
	initCache()
	val, isexists := cache[key]
	now := time.Now()
	if isexists {
		if now.After(val.time) {
			return val.data
		}
		DeleCache(key)
	}
	return nil
}

// DeleCache delete value from cache
func DeleCache(key string) {
	mu.Lock()
	defer mu.Unlock()
	initCache()
	delete(cache, key)
	loggerCache.PrintInfof("删除缓存[%s]", key)
}
