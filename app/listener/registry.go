package listener

import (
	"gin/app/event"
	"gin/pkg/serviceprovider/eventbus"
)

// All 监听器列表
func All() []eventbus.ListenerRegistration {
	return []eventbus.ListenerRegistration{
		eventbus.Bind(&UserLoginListener{}, event.UserLoginEvent{}),
		eventbus.Bind(&TestListener{}, event.UserLoginEvent{}),
	}
}
