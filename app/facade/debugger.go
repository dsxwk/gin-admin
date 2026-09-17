package facade

import (
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/eventbus"
)

// Debugger 调试器门面-调试器统一入口
func Debugger() *DebuggerFacade {
	return &DebuggerFacade{
		instance: container.Default().Get[*debugger.Debugger](serviceprovider.ServiceDebugger),
	}
}

// DebuggerFacade 调试器门面
type DebuggerFacade struct {
	instance *debugger.Debugger
}

// getInstance 获取调试器实例
func (d *DebuggerFacade) getInstance() *debugger.Debugger {
	if d == nil {
		return nil
	}
	return d.instance
}

// Start 启动调试器
func (d *DebuggerFacade) Start() {
	if inst := d.getInstance(); inst != nil {
		inst.Start()
	}
}

// Stop 停止调试器
func (d *DebuggerFacade) Stop() {
	if inst := d.getInstance(); inst != nil {
		inst.Stop()
	}
}

// IsRunning 检查调试器是否运行中
func (d *DebuggerFacade) IsRunning() bool {
	inst := d.getInstance()
	if inst == nil {
		return false
	}
	return inst.IsRunning()
}

// GetSubId 获取指定主题的订阅ID
func (d *DebuggerFacade) GetSubId(topic string) (uint64, bool) {
	inst := d.getInstance()
	if inst == nil {
		return 0, false
	}
	return inst.GetSubId(topic)
}

// GetInstance 获取原始调试器实例
func (d *DebuggerFacade) GetInstance() *debugger.Debugger {
	return d.getInstance()
}

// Bus 获取调试器使用的总线
func (d *DebuggerFacade) Bus() *eventbus.Bus {
	inst := d.getInstance()
	if inst == nil {
		return nil
	}
	return inst.Bus()
}

// Store 获取调试器使用的追踪存储
func (d *DebuggerFacade) Store() *debugger.TraceStore {
	inst := d.getInstance()
	if inst == nil {
		return nil
	}
	return inst.Store()
}
