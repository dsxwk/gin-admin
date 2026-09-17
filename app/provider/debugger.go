package provider

import (
	"context"
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/eventbus"
)

func init() {
	serviceprovider.Register(&DebuggerProvider{})
}

// DebuggerProvider 调试器服务提供者
type DebuggerProvider struct {
	instance *debugger.Debugger
}

// Name 服务提供者名称
func (p *DebuggerProvider) Name() string {
	return serviceprovider.ServiceDebugger
}

// Register 创建调试器并注册到容器
func (p *DebuggerProvider) Register(app *container.Container) {
	registry := app.Get[*eventbus.Registry](serviceprovider.ServiceEvent)
	p.instance = debugger.New(registry.Bus())
	app.Set(serviceprovider.ServiceDebugger, p.instance)
}

// Boot 启动调试器
func (p *DebuggerProvider) Boot(app *container.Container) {
	p.instance.Start()
	flag.Infof("调试器服务启动成功")
}

// Runners 后台运行任务
func (p *DebuggerProvider) Runners() []serviceprovider.Runner {
	return []serviceprovider.Runner{
		&DebuggerRunner{instance: p.instance},
	}
}

// Dependencies 依赖服务
func (p *DebuggerProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceEvent, serviceprovider.ServiceLog}
}

// DebuggerRunner 调试器后台任务
type DebuggerRunner struct {
	instance *debugger.Debugger
}

// Run 运行等待任务
func (r *DebuggerRunner) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// Stop 停止调试器
func (r *DebuggerRunner) Stop() error {
	r.instance.Stop()
	return nil
}

// Name 任务名称
func (r *DebuggerRunner) Name() string {
	return "debugger_runner"
}
