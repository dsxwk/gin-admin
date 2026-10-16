package provider

import (
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/ratelimit"
	"time"
)

// RateLimitProvider 限流服务提供者
type RateLimitProvider struct {
	manager *ratelimit.Manager
}

// Name 服务提供者名称
func (p *RateLimitProvider) Name() string {
	return serviceprovider.ServiceRateLimit
}

// Register 注册服务到容器
func (p *RateLimitProvider) Register(app *container.Container) {
	p.manager = ratelimit.NewManager(5*time.Minute, 100, 200)
	app.SetRateLimit(p.manager)
}

// Boot 启动服务
func (p *RateLimitProvider) Boot(_ *container.Container) {
	flag.Infof("限流服务启动成功")
}

// Runners 后台运行任务(用于优雅关闭)
// 返回 Runner,serviceprovider会在应用停止时自动调用Stop()
func (p *RateLimitProvider) Runners() []serviceprovider.Runner {
	manager := p.manager
	return []serviceprovider.Runner{
		serviceprovider.NewShutdownRunner("rate_limit_cleanup", func() error {
			manager.Close()
			return nil
		}),
	}
}

// Dependencies 依赖服务
func (p *RateLimitProvider) Dependencies() []string {
	return nil
}
