package event

import (
	"gin/app/facade"
	"gin/pkg/cli/base"
)

type EventListenerList struct{ base.Command }

// NewEventListenerList 创建命令
func NewEventListenerList() *EventListenerList {
	return &EventListenerList{
		Command: base.NewCommand(base.Config{
			Name:        "listener:list",
			Description: "事件监听列表",
		}),
	}
}

func (s *EventListenerList) Handle(_ map[string]string) error {
	printEventTable(facade.Event().List(), true)
	return nil
}
