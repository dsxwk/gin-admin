package facade

import (
	"gin/pkg/provider/message"
)

// Message 配置门面方法
func Message() *message.Event {
	msg := Get[*message.Event]("message")
	if msg != nil {
		return msg
	}
	return message.NewEvent()
}
