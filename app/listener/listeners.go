package listener

import (
	"gin/app/event"
	"gin/pkg/serviceprovider/eventbus"
)

// Listeners 获取监听器列表
func Listeners() []eventbus.ListenerRegistration {
	return []eventbus.ListenerRegistration{
		eventbus.Bind(&UserLoginListener{}, event.UserLoginEvent{}),
		eventbus.Bind(&TestListener{}, event.UserLoginEvent{}),
	}
}
