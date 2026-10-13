package facade

import (
	"gin/pkg/serviceprovider/ws"
)

// WS 获取WebSocket管理器
func WS() *ws.Manager {
	return App().WS()
}
