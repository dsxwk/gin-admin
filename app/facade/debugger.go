package facade

import (
	"gin/pkg/debugger"
)

// Debugger 调试器门面-调试器统一入口
var Debugger = &debuggerFacade{}

type debuggerFacade struct {
	instance *debugger.Debugger
}

// Start 启动调试器
func (d *debuggerFacade) Start() {
	if d.instance == nil {
		d.instance = debugger.NewDebugger(Message.GetBus())
	}
	d.instance.Start()
}

// Stop 停止调试器
func (d *debuggerFacade) Stop() {
	if d.instance != nil {
		d.instance.Stop()
	}
}

// IsRunning 检查调试器是否运行中
func (d *debuggerFacade) IsRunning() bool {
	if d.instance == nil {
		return false
	}
	return d.instance.IsRunning()
}

// GetInstance 获取原始调试器实例
func (d *debuggerFacade) GetInstance() *debugger.Debugger {
	return d.instance
}
