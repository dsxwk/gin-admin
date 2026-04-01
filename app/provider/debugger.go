// app/provider/debugger_provider.go
package provider

import (
	"context"
	"gin/app/facade"
	"gin/pkg/foundation"
)

func init() {
	foundation.Register(&DebuggerProvider{})
}

// DebuggerProvider 调试器服务提供者
type DebuggerProvider struct{}

func (p *DebuggerProvider) Name() string {
	return "debugger"
}

func (p *DebuggerProvider) Register(app foundation.App) {
	facade.Register("debugger", facade.Debugger)
}

func (p *DebuggerProvider) Boot(app foundation.App) {
	// 启动调试器
	facade.Debugger.Start()
	facade.Log.Info("调试器服务启动成功")
}

func (p *DebuggerProvider) Runners() []foundation.Runner {
	return []foundation.Runner{
		&DebuggerRunner{},
	}
}

func (p *DebuggerProvider) Dependencies() []string {
	return []string{"event", "log"}
}

// DebuggerRunner 调试器后台任务
type DebuggerRunner struct{}

func (r *DebuggerRunner) Run(ctx context.Context) error {
	// 等待停止信号
	<-ctx.Done()
	return nil
}

func (r *DebuggerRunner) Stop() error {
	facade.Debugger.Stop()
	return nil
}

func (r *DebuggerRunner) Name() string {
	return "debugger_runner"
}
