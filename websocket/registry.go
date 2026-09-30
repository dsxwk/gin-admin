package websocket

import "gin/pkg/serviceprovider/ws"

// All ws消息处理器列表
func All() []ws.Handler {
	return []ws.Handler{
		&MessageHandler{},
		&NotificationHandler{},
	}
}
