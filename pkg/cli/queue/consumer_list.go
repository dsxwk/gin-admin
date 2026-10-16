package queue

import (
	"fmt"
	"gin/app/facade"
	"gin/pkg/cli"
	"gin/pkg/cli/base"
	"sort"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type ConsumerList struct{ base.Command }

// NewConsumerList 创建命令
func NewConsumerList() *ConsumerList {
	return &ConsumerList{
		Command: base.NewCommand(base.Config{
			Name:        "consumer:list",
			Description: "消费者列表",
		}),
	}
}

func (s *ConsumerList) Handle(_ map[string]string) error {
	consumers := facade.Queue().Consumers()
	if len(consumers) == 0 {
		color.Yellow("no registered consumers")
		return nil

	}

	sort.Slice(consumers, func(i, j int) bool { return consumers[i].Name() < consumers[j].Name() })

	writer := cli.NewTable()
	writer.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Colors: text.Colors{text.FgGreen}},
		{Number: 2, Colors: text.Colors{text.FgYellow}},
		{Number: 3, Colors: text.Colors{text.FgWhite}},
		{Number: 4, Colors: text.Colors{text.FgWhite}},
	})
	writer.AppendHeader(table.Row{"消费者名称", "连接", "延迟队列", "描述"})

	for _, c := range consumers {
		ds := fmt.Sprintf("%v", c.IsDelay())
		writer.AppendRow(table.Row{c.Name(), c.Connection(), ds, c.Description()})
	}

	fmt.Println(writer.Render())
	color.Cyan(fmt.Sprintf("总计 %d 消费者", len(consumers)))
	return nil
}
