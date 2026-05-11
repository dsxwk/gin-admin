package provider

import (
	"gin/app/facade"
	"gin/pkg/foundation"
	"gin/pkg/provider/cache"
)

func init() {
	foundation.Register(&CacheProvider{})
}

// CacheProvider 缓存服务提供者
type CacheProvider struct{}

// Name 服务提供者名称
func (p *CacheProvider) Name() string {
	return "cache"
}

// Register 注册服务到门面
func (p *CacheProvider) Register(app foundation.App) {
	conf := facade.Config()
	// 注册默认缓存
	facade.Register[*cache.CacheProxy]("cache", cache.NewCache(conf))
	// 注册Redis缓存
	facade.Register[*cache.CacheProxy]("redis_cache", cache.NewRedisCache(conf))
	// 注册内存缓存
	facade.Register[*cache.CacheProxy]("memory_cache", cache.NewMemoryCache(conf))
	// 注册磁盘缓存
	facade.Register[*cache.CacheProxy]("disk_cache", cache.NewDiskCache(conf))
}

// Boot 启动服务
func (p *CacheProvider) Boot(app foundation.App) {
	facade.Log().Info("缓存服务启动成功")
}

// Dependencies 依赖服务
func (p *CacheProvider) Dependencies() []string {
	return []string{"config"}
}
