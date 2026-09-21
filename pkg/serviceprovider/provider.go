package serviceprovider

import (
	"context"
	"gin/pkg/container"
)

// ServiceProvider 服务提供者
type ServiceProvider interface {
	// Name 服务提供者名称
	Name() string
	// Register 注册服务到容器
	Register(app *container.Container)
	// Boot 启动服务
	Boot(app *container.Container)
}

// Dependencies 依赖关系
type Dependencies interface {
	ServiceProvider
	// Dependencies 依赖其他服务提供者名称列表
	Dependencies() []string
}

// Runners 后台运行任务
type Runners interface {
	ServiceProvider
	// Runners 后台运行任务
	Runners() []Runner
}

// Runner 后台运行任务接口
type Runner interface {
	// Run 运行任务
	Run(ctx context.Context) error
	// Stop 停止任务
	Stop() error
	// Name 任务名称
	Name() string
}
