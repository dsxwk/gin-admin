package facade

import (
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/cache"

	"github.com/go-redis/redis/v8"
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
	manager := container.Default().Get[*cache.Manager](serviceprovider.ServiceCache)
	if manager == nil {
		return nil
	}
	return manager.Cache(cacheType...)
}

// Redis 获取Redis缓存实例
func Redis() *cache.RedisCache {
	manager := container.Default().Get[*cache.Manager](serviceprovider.ServiceCache)
	if manager == nil {
		return nil
	}
	return manager.Redis()
}

// RedisClient 获取Redis客户端
func RedisClient() *redis.Client {
	instance := Redis()
	if instance == nil {
		return nil
	}
	return instance.Client()
}
