package {{.Package}}

import (
    "fmt"
    "gin/app/event"
    "gin/common/base"
    "gin/pkg/eventbus"
    "github.com/goccy/go-json"
    "time"
)

type {{.Name}}Listener struct{}

func (l *{{.Name}}Listener) Handle(e base.Event) {
    ev, ok := e.(event.{{.EventName}})
	if !ok {
		return
	}

	data, _ := json.Marshal(e)
    fmt.Printf("收到事件: %s 事件描述: %s 事件数据: %s, 时间: %s\n", ev.Name(), ev.Description(), data, time.Now().Format("2006-01-02 15:04:05"))
}

func init() {
	eventbus.Register(&{{.Name}}Listener{}, event.{{.EventName}}{})
}
