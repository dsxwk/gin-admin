package eventbus

// Listener 泛型监听器接口,处理指定事件类型
type Listener[T Event] interface {
	Handle(e T) // 处理事件
}
