package facade

import (
	"gin/pkg/provider/debugger"
)

// Debugger 调试器门面-调试器统一入口
func Debugger() *DebuggerFacade {
	return &DebuggerFacade{}
}

type DebuggerFacade struct{}

// instance 从Manager获取调试器实例
func (d *DebuggerFacade) instance() *debugger.Debugger {
	dbg := Get[*debugger.Debugger]("debugger")
	if dbg != nil {
		return dbg
	}
	return debugger.NewDebugger(Message())
}

// Start 启动调试器
func (d *DebuggerFacade) Start() {
	if inst := d.instance(); inst != nil {
		inst.Start()
	}
}

// Stop 停止调试器
func (d *DebuggerFacade) Stop() {
	if inst := d.instance(); inst != nil {
		inst.Stop()
	}
}

// IsRunning 检查调试器是否运行中
func (d *DebuggerFacade) IsRunning() bool {
	inst := d.instance()
	if inst == nil {
		return false
	}
	return inst.IsRunning()
}

// GetInstance 获取原始调试器实例
func (d *DebuggerFacade) GetInstance() *debugger.Debugger {
	return d.instance()
}
