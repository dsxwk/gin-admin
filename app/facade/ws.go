package facade

import (
	"gin/pkg/container"
	"gin/pkg/serviceprovider/ws"
)

// WS 获取WebSocket管理器
func WS() *ws.Manager {
	return container.Default().WS()
}
