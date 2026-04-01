package event

import (
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg/cli"
	"github.com/fatih/color"
)

type EventList struct{}

func (s *EventList) Name() string {
	return "event:list"
}

func (s *EventList) Description() string {
	return "事件列表"
}

func (s *EventList) Help() []base.CommandOption {
	return []base.CommandOption{}
}

func (s *EventList) Execute(args []string) {
	list := facade.Event.List()
	for _, v := range list {
		color.Green(fmt.Sprintf("%-8s %-35s", v.Name, v.Description))
	}
}

func init() {
	cli.Register(&EventList{})
}
