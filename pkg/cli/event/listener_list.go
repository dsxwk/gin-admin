package event

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg/cli"
	"gin/pkg/eventbus"
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

func (s *EventListenerList) Execute(args []string) {
	facade.Event[eventbus.Event]().Debug()
}

func init() {
	cli.Register(&EventListenerList{})
}
