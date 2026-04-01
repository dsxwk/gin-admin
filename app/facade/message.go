package facade

import (
	"gin/pkg/message"
)

// Message 消息事件门面-事件总线统一入口
// 使用示例:
//
//	facade.Message.Publish("user.login", eventData)
//	facade.Message.Subscribe("user.login", func(data any) { ... })
var Message = &messageFacade{}

type messageFacade struct{}

// Publish 发布事件
func (e *messageFacade) Publish(topic string, event any) {
	message.GetEventBus().Publish(topic, event)
}

// Subscribe 同步订阅事件
func (e *messageFacade) Subscribe(topic string, fn func(any)) uint64 {
	return message.GetEventBus().Subscribe(topic, fn)
}

// SubscribeAsync 异步订阅事件
func (e *messageFacade) SubscribeAsync(topic string, fn func(any)) uint64 {
	return message.GetEventBus().SubscribeAsync(topic, fn)
}

// Unsubscribe 取消订阅
func (e *messageFacade) Unsubscribe(topic string, id uint64) bool {
	return message.GetEventBus().Unsubscribe(topic, id)
}

// GetBus 获取原始事件总线
func (e *messageFacade) GetBus() *message.EventBus {
	return message.GetEventBus()
}
