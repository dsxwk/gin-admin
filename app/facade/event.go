package facade

import (
	"context"
	"gin/pkg/serviceprovider/eventbus"
)

// Event 门面函数
// 使用示例:
//
//	facade.Event().Register[event.UserLoginEvent](listener, event)
//	facade.Event().Publish[event.UserLoginEvent](ctx, event)
func Event() EventFacade {
	return EventFacade{}
}

type EventFacade struct{}

// Register 注册监听器
func (f EventFacade) Register[T eventbus.Event](listener eventbus.Listener[T], e T) {
	eventbus.Register[T](listener, e)
}

// Publish 发布事件
func (f EventFacade) Publish[T eventbus.Event](ctx context.Context, e T) {
	eventbus.Publish[T](ctx, e)
}

// List 获取事件列表
func (f EventFacade) List() []eventbus.EventInfo {
	return eventbus.EventList()
}

// Debug 打印事件
func (f EventFacade) Debug() {
	eventbus.DebugPrint()
}
