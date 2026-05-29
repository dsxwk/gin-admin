package provider

import (
	"gin/app/facade"
	"gin/common/flag"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/cache"
)

func init() {
	serviceprovider.Register(&CacheProvider{})
}

// CacheProvider 缓存服务提供者
type CacheProvider struct{}

// Name 服务提供者名称
func (p *CacheProvider) Name() string {
	return "cache"
}

// Register 注册服务到门面
func (p *CacheProvider) Register(app serviceprovider.App) {
	cfg := facade.Config()
	// 注册默认缓存
	facade.Register[*cache.CacheProxy](cfg.Cache.Driver, cache.NewCache(cfg.Cache.Driver, cfg))
}

// Boot 启动服务
func (p *CacheProvider) Boot(app serviceprovider.App) {
	flag.Infof("缓存服务启动成功")
}

// Dependencies 依赖服务
func (p *CacheProvider) Dependencies() []string {
	return []string{"config"}
}
