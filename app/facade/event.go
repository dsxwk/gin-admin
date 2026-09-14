package facade

import (
	"context"
	"gin/pkg/serviceprovider/eventbus"
	"sync"
)

var (
	eventOnce     sync.Once
	eventRegistry *eventbus.Registry
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

// Registry 获取业务事件注册表
func (f EventFacade) Registry() *eventbus.Registry {
	eventOnce.Do(func() {
		eventRegistry = Get[*eventbus.Registry]("event")
		if eventRegistry != nil {
			return
		}

		eventRegistry = eventbus.NewRegistry(eventbus.NewBus())
		Register[*eventbus.Registry]("event", eventRegistry)
	})
	return eventRegistry
}

// Bus 获取底层事件总线
func (f EventFacade) Bus() *eventbus.Bus {
	registry := f.Registry()
	if registry != nil {
		return registry.Bus()
	}
	return eventbus.NewBus()
}

// Register 注册监听器
func (f EventFacade) Register[T eventbus.Event](listener eventbus.Listener[T], e T) {
	registry := f.Registry()
	if registry == nil {
		return
	}
	registry.Register(listener, e)
}

// Publish 发布事件
func (f EventFacade) Publish[T eventbus.Event](ctx context.Context, e T) {
	registry := f.Registry()
	if registry == nil {
		return
	}
	registry.Publish(ctx, e)
}

// List 获取事件列表
func (f EventFacade) List() []eventbus.EventInfo {
	registry := f.Registry()
	if registry == nil {
		return nil
	}
	return registry.EventList()
}
