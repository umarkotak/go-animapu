package app_cache

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

var (
	appCache *cache.Cache
)

func InitAppCache() {
	appCache = cache.New(1*time.Minute, 3*time.Minute)
}

func GetAppCache() *cache.Cache {
	return appCache
}
