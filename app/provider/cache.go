package provider

import (
	"gin/common/flag"
	"gin/config"
	"gin/pkg/container"
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
	return serviceprovider.ServiceCache
}

// Register 注册服务到容器
func (p *CacheProvider) Register(app *container.Container) {
	cfg := app.Get[*config.Config](serviceprovider.ServiceConfig)
	app.Set(serviceprovider.ServiceCache, cache.NewManager(cfg))
}

// Boot 启动服务
func (p *CacheProvider) Boot(app *container.Container) {
	flag.Infof("缓存服务启动成功")
}

// Dependencies 依赖服务
func (p *CacheProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig}
}
