package provider

import (
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/cache"
)

// CacheProvider 缓存服务提供者
type CacheProvider struct{}

// Name 服务提供者名称
func (p *CacheProvider) Name() string {
	return serviceprovider.ServiceCache
}

// Register 注册服务到容器
func (p *CacheProvider) Register(app *container.Container) {
	app.SetCache(cache.NewManager(app.Config()))
}

// Boot 启动服务
func (p *CacheProvider) Boot(_ *container.Container) {
	flag.Infof("缓存服务启动成功")
}

// Dependencies 依赖服务
func (p *CacheProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig}
}
