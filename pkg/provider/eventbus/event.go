package eventbus

// Event 事件接口
type Event interface {
	Name() string        // 事件名称
	Description() string // 事件描述
}
