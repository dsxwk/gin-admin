package ws

import "gin/pkg/errcode"

// WsErrCodePrefix WebSocket错误码前缀
const WsErrCodePrefix = 800

var (
	// ErrClosed WebSocket连接已关闭
	ErrClosed = errcode.NewError(1, "websocket连接已关闭").WithPrefix(WsErrCodePrefix)
	// ErrTooManyConnections WebSocket连接数已达到上限
	ErrTooManyConnections = errcode.NewError(2, "websocket连接数已达到上限").WithPrefix(WsErrCodePrefix)
	// ErrSlowClient WebSocket客户端发送过慢
	ErrSlowClient = errcode.NewError(3, "websocket客户端发送过慢").WithPrefix(WsErrCodePrefix)
)
