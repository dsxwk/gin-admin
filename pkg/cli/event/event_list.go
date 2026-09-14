package event

import (
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg/cli"
	"gin/pkg/serviceprovider/eventbus"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
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

func (s *EventList) Execute(values map[string]string) {
	printEventTable(facade.Event().List(), false)
}

// printEventTable 打印事件表格
func printEventTable(events []eventbus.EventInfo, showListeners bool) {
	if len(events) == 0 {
		color.Yellow("暂无注册的事件")
		return
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Name < events[j].Name
	})

	nameWidth := max(20, runewidth.StringWidth("事件名称"))
	descWidth := max(35, runewidth.StringWidth("描述"))
	for _, event := range events {
		nameWidth = max(nameWidth, runewidth.StringWidth(event.Name))
		descWidth = max(descWidth, runewidth.StringWidth(event.Description))
	}

	totalWidth := nameWidth + descWidth + 7

	color.Yellow("┌" + strings.Repeat("─", totalWidth-2) + "┐")
	color.Yellow(fmt.Sprintf(
		"│ %s   %s %s",
		color.HiWhiteString(padRight("事件名称", nameWidth)),
		color.HiWhiteString(padRight("描述", descWidth)),
		color.YellowString("│"),
	))
	color.Yellow("├" + strings.Repeat("─", totalWidth-2) + "┤")

	listeners := 0
	for _, event := range events {
		color.Yellow(fmt.Sprintf(
			"│ %s   %s %s",
			color.GreenString(padRight(event.Name, nameWidth)),
			color.WhiteString(padRight(event.Description, descWidth)),
			color.YellowString("│"),
		))

		if !showListeners {
			continue
		}

		for index, listener := range event.Listeners {
			prefix := "├─ "
			if index == len(event.Listeners)-1 {
				prefix = "└─ "
			}
			listeners++

			row := strings.Repeat(" ", nameWidth+2) + prefix + listener
			color.Yellow(fmt.Sprintf("│%-*s│", totalWidth-2, row))
		}
	}

	color.Yellow("└" + strings.Repeat("─", totalWidth-2) + "┘")

	if showListeners {
		color.Cyan(fmt.Sprintf("总计 %d 个事件 %d 个监听\n", len(events), listeners))
		return
	}

	color.Cyan(fmt.Sprintf("总计 %d 个事件\n", len(events)))
}

// padRight 右侧填充空格,支持中文字符
func padRight(s string, width int) string {
	currentWidth := runewidth.StringWidth(s)
	if currentWidth >= width {
		return s
	}
	padding := width - currentWidth
	return s + strings.Repeat(" ", padding)
}

func init() {
	cli.Register(&EventList{})
}
