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

type ProducerList struct{ base.Command }

// NewProducerList 创建命令
func NewProducerList() *ProducerList {
	return &ProducerList{
		Command: base.NewCommand(base.Config{
			Name:        "producer:list",
			Description: "生产者列表",
		}),
	}
}

func (s *ProducerList) Handle(_ map[string]string) error {
	producers := facade.Queue().Producers()
	if len(producers) == 0 {
		color.Yellow("no registered producers")
		return nil

	}

	sort.Slice(producers, func(i, j int) bool { return producers[i].Name() < producers[j].Name() })

	writer := cli.NewTable()
	writer.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Colors: text.Colors{text.FgGreen}},
		{Number: 2, Colors: text.Colors{text.FgYellow}},
		{Number: 3, Colors: text.Colors{text.FgWhite}},
		{Number: 4, Colors: text.Colors{text.FgWhite}},
		{Number: 5, Colors: text.Colors{text.FgWhite}},
	})
	writer.AppendHeader(table.Row{"生产者名称", "连接", "延迟队列", "延迟(ms)", "描述"})

	for _, p := range producers {
		ds := fmt.Sprintf("%v", p.IsDelay())
		dms := fmt.Sprintf("%dms", p.DelayMs())
		writer.AppendRow(table.Row{p.Name(), p.Connection(), ds, dms, p.Description()})
	}

	fmt.Println(writer.Render())
	color.Cyan(fmt.Sprintf("总计 %d 生产者", len(producers)))
	return nil
}
