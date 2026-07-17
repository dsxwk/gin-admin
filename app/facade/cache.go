package facade

import (
	"gin/pkg/serviceprovider/cache"
)

// Cache 获取缓存实例
// 不传参数返回默认缓存，传参返回指定缓存
//
// 使用示例:
//
//	// 默认缓存
//	facade.Cache().Set("key", "value", 5*time.Minute)
//
//	// 指定缓存
//	facade.Cache("redis").Set("key", "value", 5*time.Minute)
//	facade.Cache("memory").Get("key")
func Cache(cacheType ...string) *cache.CacheProxy {
	name := Config().Cache.Driver
	if len(cacheType) > 0 && cacheType[0] != "" {
		name = cacheType[0]
	}

	_cache := Get[*cache.CacheProxy](name)
	if _cache != nil {
		return _cache
	}
	cp := cache.NewCache(name, Config())
	Register[*cache.CacheProxy](name, cp)
	return cp
}

// ResetCache 重置缓存
func ResetCache(cacheType ...string) *cache.CacheProxy {
	name := Config().Cache.Driver
	if len(cacheType) > 0 && cacheType[0] != "" {
		name = cacheType[0]
	}
	cp := cache.ResetCache(name, Config())
	Register[*cache.CacheProxy](name, cp)
	return cp
}
