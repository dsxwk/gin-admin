package facade

import (
	"context"
	"gin/pkg/cache"
)

// Cache 缓存门面-缓存统一入口
// 使用示例:
//
//	facade.Cache.Set("key", "value", 5*time.Minute)
//	value, ok := facade.Cache.Store("redis").Get("key")
//	cache := facade.Cache.WithContext(ctx)  // 获取带上下文的缓存
var Cache = &cacheFacade{}

type cacheFacade struct{}

// Store 获取指定类型的缓存存储
// 参数说明:
//   - cacheType: 缓存类型，可选 "redis", "memory", "disk"，不传时返回默认缓存
//
// 使用示例:
//
//	facade.Cache.Store()           // 获取默认缓存
//	facade.Cache.Store("redis")    // 获取Redis缓存
//	facade.Cache.Store("memory")   // 获取内存缓存
func (c *cacheFacade) Store(cacheType ...string) *cache.CacheProxy {
	typ := "default"
	if len(cacheType) > 0 && cacheType[0] != "" {
		typ = cacheType[0]
	}

	var key string
	switch typ {
	case "redis":
		key = "redis_cache"
	case "memory":
		key = "memory_cache"
	case "disk":
		key = "disk_cache"
	default:
		key = "cache"
	}

	if _cache := Get(key); _cache != nil {
		return _cache.(*cache.CacheProxy)
	}
	return cache.NewCache(Config.Get())
}

// WithContext 获取带上下文的缓存
func (c *cacheFacade) WithContext(ctx context.Context) *cache.CacheProxy {
	return c.Store().WithContext(ctx)
}

// Redis 获取Redis缓存
func (c *cacheFacade) Redis() *cache.CacheProxy {
	return c.Store("redis")
}

// Memory 获取内存缓存
func (c *cacheFacade) Memory() *cache.CacheProxy {
	return c.Store("memory")
}

// Disk 获取磁盘缓存
func (c *cacheFacade) Disk() *cache.CacheProxy {
	return c.Store("disk")
}
