package event

import (
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg"
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
	if len(list) == 0 {
		color.Yellow("暂无注册的事件")
		return
	}
	// 计算最大宽度
	maxNameLen := 0
	maxDescLen := 0
	for _, v := range list {
		if len(v.Name) > maxNameLen {
			maxNameLen = len(v.Name)
		}
		if len(v.Description) > maxDescLen {
			maxDescLen = len(v.Description)
		}
	}

	// 设置最小宽度
	if maxNameLen < 10 {
		maxNameLen = 10
	}
	if maxDescLen < 20 {
		maxDescLen = 20
	}
	fmt.Println(pkg.Sprintf("%-*s %-*s", maxNameLen+2, "Name", maxDescLen+2, "Description"))
	for _, v := range list {
		fmt.Printf("%s %s",
			color.GreenString(fmt.Sprintf("%-*s", maxNameLen+2, v.Name)),
			color.YellowString(fmt.Sprintf("%-*s", maxDescLen+2, v.Description)),
		)
	}
}

func init() {
	cli.Register(&EventList{})
}
