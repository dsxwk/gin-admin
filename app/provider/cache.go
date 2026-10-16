package provider

import (
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/cache"
)

// CacheProvider 缓存服务提供者
type CacheProvider struct {
	manager *cache.Manager
}

// Name 服务提供者名称
func (p *CacheProvider) Name() string {
	return serviceprovider.ServiceCache
}

// Register 注册服务到容器
func (p *CacheProvider) Register(app *container.Container) {
	p.manager = cache.NewManager(app.Config())
	app.SetCache(p.manager)
}

// Boot 启动服务
func (p *CacheProvider) Boot(_ *container.Container) {
	flag.Infof("缓存服务启动成功")
}

// Dependencies 依赖服务
func (p *CacheProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig}
}

// Runners 后台运行任务
func (p *CacheProvider) Runners() []serviceprovider.Runner {
	if p.manager == nil {
		return nil
	}

	return []serviceprovider.Runner{
		serviceprovider.NewShutdownRunner("cache_shutdown", p.manager.Close),
	}
}
