package cache

import (
	"QuickStone/src/config"

	"github.com/coocood/freecache"
)

var DefaultExpireSeconds int

var cache *freecache.Cache

func init() {
	cacheSize := config.FreeCacheSize
	DefaultExpireSeconds = config.FreeCacheDefaultExpireSeconds
	cache = freecache.NewCache(cacheSize)
}

func Set(obj any) error {
	return nil
}
