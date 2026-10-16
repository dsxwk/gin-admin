package provider

import (
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/terminal"
)

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
	p.instance = debugger.New(app.Event().Bus())
	app.SetDebugger(p.instance)
}

// Boot 启动调试器
func (p *DebuggerProvider) Boot(_ *container.Container) {
	p.instance.Start()
	terminal.Infof("调试器服务启动成功")
}

// Runners 后台运行任务
func (p *DebuggerProvider) Runners() []serviceprovider.Runner {
	return []serviceprovider.Runner{
		serviceprovider.NewShutdownRunner("debugger_runner", func() error {
			p.instance.Stop()
			return nil
		}),
	}
}

// Dependencies 依赖服务
func (p *DebuggerProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceEvent, serviceprovider.ServiceLog}
}
