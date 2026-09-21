package listener

import (
	"gin/app/event"
	"gin/pkg/serviceprovider/eventbus"
)

type listenerEntry func(*eventbus.Registry)

// Register 注册业务事件监听器
func Register(registry *eventbus.Registry) {
	for _, entry := range listenerEntries() {
		entry(registry)
	}
}

// listenerEntries 获取监听器注册列表
func listenerEntries() []listenerEntry {
	return []listenerEntry{
		listenerRegister(&UserLoginListener{}, event.UserLoginEvent{}),
		listenerRegister(&TestListener{}, event.UserLoginEvent{}),
	}
}

// listenerRegister 创建监听器注册函数
func listenerRegister[T eventbus.Event](listener eventbus.Listener[T], e T) listenerEntry {
	return func(registry *eventbus.Registry) {
		registry.Register(listener, e)
	}
}
