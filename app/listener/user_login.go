package listener

import (
	"fmt"
	"gin/app/event"
	"gin/app/facade"
	"time"
)

type UserLoginListener struct{}

func (l *UserLoginListener) Handle(e event.UserLoginEvent) {
	fmt.Printf(
		"收到事件: %s 描述: %s 数据: %T 时间: %s\n",
		e.Name(),
		e.Description(),
		e,
		time.Now().Format("2006-01-02 15:04:05"),
	)
}

func init() {
	facade.Event().Register(&UserLoginListener{}, event.UserLoginEvent{})
}
