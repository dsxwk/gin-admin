package base

// Listener 泛型接口,处理指定事件类型
//type Listener[T Event] interface {
//	Handle(e T) // 处理事件
//}

// Listener 监听器接口
type Listener interface {
	Handle(event Event)
}
