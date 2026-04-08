package facade

import (
	"context"
	"gin/pkg/eventbus"
)

// Event 门面
var Event = &eventFacade{}

type eventFacade struct{}

// Register 注册监听器
func (f *eventFacade) Register(listener eventbus.Listener, event eventbus.Event) {
	eventbus.Register(listener, event)
}

// Publish 发布事件
func (f *eventFacade) Publish(ctx context.Context, e eventbus.Event) {
	eventbus.Publish(ctx, e)
}

// List 获取事件列表
func (f *eventFacade) List() []eventbus.EventInfo {
	return eventbus.EventList()
}

// Debug 打印事件
func (f *eventFacade) Debug() {
	eventbus.DebugPrint()
}
