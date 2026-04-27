package {{.Package}}

import (
    "fmt"
    "gin/app/event"
    "gin/pkg/eventbus"
    "time"
)

type {{.Name}}Listener struct{}

func (l *{{.Name}}Listener) Handle(e event.{{.EventName}}) {
    fmt.Printf(
        "收到事件: %s 描述: %s 数据: %T 时间: %s\n",
        e.Name(),
        e.Description(),
        e,
        time.Now().Format("2006-01-02 15:04:05"),
    )
}

func init() {
	eventbus.Register(&{{.Name}}Listener{}, event.{{.EventName}}{})
}
