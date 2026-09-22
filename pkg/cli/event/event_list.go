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
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

func init() {
	cli.Register(&EventList{})
}
