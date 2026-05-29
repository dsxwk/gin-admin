package provider

import (
	"context"
	"gin/app/facade"
	"gin/common/flag"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/debugger"
)

func init() {
	serviceprovider.Register(&DebuggerProvider{})
}

// DebuggerProvider 调试器服务提供者
type DebuggerProvider struct{}

func (p *DebuggerProvider) Name() string {
	return "debugger"
}

func (p *DebuggerProvider) Register(app serviceprovider.App) {
	facade.Register[*debugger.Debugger]("debugger", debugger.NewDebugger(facade.Message()))
}

func (p *DebuggerProvider) Boot(app serviceprovider.App) {
	// 启动调试器
	facade.Debugger().Start()
	flag.Infof("调试器服务启动成功")
}

func (p *DebuggerProvider) Runners() []serviceprovider.Runner {
	return []serviceprovider.Runner{
		&DebuggerRunner{},
	}
}

func (p *DebuggerProvider) Dependencies() []string {
	return []string{"message", "event", "log"}
}

// DebuggerRunner 调试器后台任务
type DebuggerRunner struct{}

func (r *DebuggerRunner) Run(ctx context.Context) error {
	// 等待停止信号
	<-ctx.Done()
	return nil
}

func (r *DebuggerRunner) Stop() error {
	facade.Debugger().Stop()
	return nil
}

func (r *DebuggerRunner) Name() string {
	return "debugger_runner"
}
