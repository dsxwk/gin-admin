package facade

import (
	"gin/pkg/serviceprovider/debugger"
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

// GetSubId 获取指定主题的订阅ID
func (d *DebuggerFacade) GetSubId(topic string) (uint64, bool) {
	inst := d.instance()
	if inst == nil {
		return 0, false
	}
	return inst.GetSubId(topic)
}

// GetInstance 获取原始调试器实例
func (d *DebuggerFacade) GetInstance() *debugger.Debugger {
	return d.instance()
}
