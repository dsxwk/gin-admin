package eventbus

import "context"

// Event 业务事件接口
type Event interface {
	Name() string
	Description() string
}

// Listener 泛型事件监听器
type Listener[T Event] interface {
	Handle(event T)
}

// PublishedEvent 已发布业务事件
type PublishedEvent struct {
	Context     context.Context
	Name        string
	Description string
	Data        Event
}
