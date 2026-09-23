package provider

import (
	"context"
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
	return []serviceprovider.Runner{
		&RateLimitCleanupRunner{manager: p.manager},
	}
}

// Dependencies 依赖服务
func (p *RateLimitProvider) Dependencies() []string {
	return nil
}

// RateLimitCleanupRunner 限流清理任务
type RateLimitCleanupRunner struct {
	manager *ratelimit.Manager
}

// Run 运行清理任务
func (r *RateLimitCleanupRunner) Run(ctx context.Context) error {
	// 等待停止信号
	<-ctx.Done()
	return nil
}

// Stop 停止时关闭限流器
func (r *RateLimitCleanupRunner) Stop() error {
	r.manager.Close()
	return nil
}

// Name 任务名称
func (r *RateLimitCleanupRunner) Name() string {
	return "rate_limit_cleanup"
}
