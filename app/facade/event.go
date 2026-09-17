package facade

import (
	"context"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/eventbus"
)

// Event 门面函数
// 使用示例:
//
//	facade.Event().Register[event.UserLoginEvent](listener, event)
//	facade.Event().Publish[event.UserLoginEvent](ctx, event)
func Event() EventFacade {
	registry := container.Default().Get[*eventbus.Registry](serviceprovider.ServiceEvent)
	if registry == nil {
		registry = eventbus.DefaultRegistry()
	}

	return EventFacade{
		registry: registry,
	}
}

// EventFacade 事件门面
type EventFacade struct {
	registry *eventbus.Registry
}

// Registry 获取业务事件注册表
func (f EventFacade) Registry() *eventbus.Registry {
	return f.registry
}

// Bus 获取底层事件总线
func (f EventFacade) Bus() *eventbus.Bus {
	if f.registry == nil {
		return nil
	}
	return f.registry.Bus()
}

// Register 注册监听器
func (f EventFacade) Register[T eventbus.Event](listener eventbus.Listener[T], e T) {
	if f.registry == nil {
		return
	}
	f.registry.Register(listener, e)
}

// Publish 发布事件
func (f EventFacade) Publish[T eventbus.Event](ctx context.Context, e T) {
	if f.registry == nil {
		return
	}
	f.registry.Publish(ctx, e)
}

// List 获取事件列表
func (f EventFacade) List() []eventbus.EventInfo {
	if f.registry == nil {
		return nil
	}
	return f.registry.EventList()
}
