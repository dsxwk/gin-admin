package facade

import (
	"context"
	"gin/pkg/eventbus"
)

// Event 门面泛型函数
// 使用示例:
//
//	facade.Event[event.UserLoginEvent]().Register(listener, event)
//	facade.Event[event.UserLoginEvent]().Publish(ctx, event)
func Event[T eventbus.Event]() EventFacade[T] {
	return EventFacade[T]{}
}

type EventFacade[T eventbus.Event] struct{}

// Register 注册监听器
func (f EventFacade[T]) Register(listener eventbus.Listener[T], e T) {
	eventbus.Register[T](listener, e)
}

// Publish 发布事件
func (f EventFacade[T]) Publish(ctx context.Context, e T) {
	eventbus.Publish[T](ctx, e)
}

// List 获取事件列表
func (f EventFacade[T]) List() []eventbus.EventInfo {
	return eventbus.EventList()
}

// Debug 打印事件
func (f EventFacade[T]) Debug() {
	eventbus.DebugPrint()
}
