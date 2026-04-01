package facade

import (
	"context"
	"gin/common/base"
	"gin/pkg/eventbus"
)

// Event 门面
var Event = &eventFacade{}

type eventFacade struct{}

// Register 注册监听器
func (f *eventFacade) Register(listener base.Listener, event base.Event) {
	eventbus.Register(listener, event)
}

// Publish 发布事件
func (f *eventFacade) Publish(ctx context.Context, e base.Event) {
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
