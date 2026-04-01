package provider

import (
	"context"
	"gin/app/facade"
	"gin/pkg/foundation"
	"os"
)

func init() {
	foundation.Register(&RateLimitProvider{})
}

// RateLimitProvider 限流服务提供者
type RateLimitProvider struct{}

// Name 服务提供者名称
func (p *RateLimitProvider) Name() string {
	return "rate_limit"
}

// Register 注册服务到门面
func (p *RateLimitProvider) Register(app foundation.App) {
	facade.Register("rate_limit", facade.RateLimiter)
}

// Boot 启动服务
func (p *RateLimitProvider) Boot(app foundation.App) {
	// 初始化限流
	facade.RateLimiter.Init()
	facade.Log.Info("限流服务启动成功")
}

// Runners 后台运行任务(用于优雅关闭)
// 返回 Runner,foundation会在应用停止时自动调用Stop()
func (p *RateLimitProvider) Runners() []foundation.Runner {
	return []foundation.Runner{
		&RateLimitCleanupRunner{},
	}
}

// Dependencies 依赖服务
func (p *RateLimitProvider) Dependencies() []string {
	return []string{"config", "log"}
}

// RateLimitCleanupRunner 限流清理任务
type RateLimitCleanupRunner struct{}

// Run 运行清理任务
func (r *RateLimitCleanupRunner) Run(ctx context.Context) error {
	// 等待停止信号
	<-ctx.Done()
	return nil
}

// Stop 停止时关闭限流器
func (r *RateLimitCleanupRunner) Stop() error {
	facade.RateLimiter.Shutdown()
	if os.Getenv("CLI_MODE") != "true" {
		facade.Log.Info("限流服务已关闭")
	}
	return nil
}

// Name 任务名称
func (r *RateLimitCleanupRunner) Name() string {
	return "rate_limit_cleanup"
}
