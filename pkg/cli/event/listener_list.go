package event

import (
	"gin/app/facade"
	"gin/common/base"
)

type EventListenerList struct{}

func (s *EventListenerList) Name() string {
	return "listener:list"
}

func (s *EventListenerList) Description() string {
	return "事件监听列表"
}

func (s *EventListenerList) Help() []base.CommandOption {
	return []base.CommandOption{}
}

func (s *EventListenerList) Execute(_ map[string]string) {
	printEventTable(facade.Event().List(), true)
}
