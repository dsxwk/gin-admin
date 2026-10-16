package event

import (
	"fmt"
	"gin/app/facade"
	"gin/pkg/cli"
	"gin/pkg/cli/base"
	"gin/pkg/serviceprovider/eventbus"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type EventList struct{ base.Command }

// NewEventList 创建命令
func NewEventList() *EventList {
	return &EventList{
		Command: base.NewCommand(base.Config{
			Name:        "event:list",
			Description: "事件列表",
		}),
	}
}

func (s *EventList) Handle(_ map[string]string) error {
	printEventTable(facade.Event().List(), false)
	return nil
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

	writer := cli.NewTable()
	writer.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Colors: text.Colors{text.FgGreen}},
		{Number: 2, Colors: text.Colors{text.FgWhite}},
		{Number: 3, Colors: text.Colors{text.FgCyan}},
	})
	if showListeners {
		writer.AppendHeader(table.Row{"事件名称", "描述", "监听器"})
	} else {
		writer.AppendHeader(table.Row{"事件名称", "描述"})
	}

	listeners := 0
	for _, event := range events {
		if !showListeners {
			writer.AppendRow(table.Row{event.Name, event.Description})
			continue
		}

		listeners += len(event.Listeners)
		writer.AppendRow(table.Row{
			event.Name,
			event.Description,
			strings.Join(event.Listeners, "\n"),
		})
	}

	fmt.Println(writer.Render())

	if showListeners {
		color.Cyan(fmt.Sprintf("总计 %d 个事件 %d 个监听\n", len(events), listeners))
		return
	}

	color.Cyan(fmt.Sprintf("总计 %d 个事件\n", len(events)))
}
