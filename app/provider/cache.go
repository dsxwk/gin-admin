package provider

import (
	"context"
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
		&cacheShutdownRunner{manager: p.manager},
	}
}

// cacheShutdownRunner 缓存关闭任务
type cacheShutdownRunner struct {
	manager *cache.Manager
}

// Run 等待停止信号
func (r *cacheShutdownRunner) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// Stop 关闭缓存驱动
func (r *cacheShutdownRunner) Stop() error {
	return r.manager.Close()
}

// Name 任务名称
func (r *cacheShutdownRunner) Name() string {
	return "cache_shutdown"
}
