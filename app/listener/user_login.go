package listener

import (
	"fmt"
	"gin/app/event"
	"gin/pkg/eventbus"
	"github.com/goccy/go-json"
	"time"
)

type UserLoginListener struct{}

func (l *UserLoginListener) Handle(e eventbus.Event) {
	ev, ok := e.(event.UserLoginEvent)
	if !ok {
		return
	}

	data, _ := json.Marshal(ev)

	fmt.Printf(
		"收到事件: %s 描述: %s 数据: %s 时间: %s\n",
		ev.Name(),
		ev.Description(),
		data,
		time.Now().Format("2006-01-02 15:04:05"),
	)
}

func init() {
	eventbus.Register(&UserLoginListener{}, event.UserLoginEvent{})
}
