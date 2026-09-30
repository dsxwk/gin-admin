package ws

import "context"

// Connection WebSocket连接
type Connection interface {
	ID() string                 // 客户端ID
	UserID() int64              // 用户ID
	BindUser(userID int64)      // 绑定用户
	Send(message Message) error // 发送消息
	Close()                     // 关闭连接
}

// Handler WebSocket消息处理器
type Handler interface {
	Handle(ctx context.Context, connection Connection, message Message) error // 处理消息
}

// HandlerMatcher 消息匹配处理器
type HandlerMatcher interface {
	Match(message Message) bool // 匹配消息
}
