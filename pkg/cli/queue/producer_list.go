package queue

import (
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg/cli"
	"sort"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type ProducerList struct{}

func (s *ProducerList) Name() string               { return "producer:list" }
func (s *ProducerList) Description() string        { return "生产者列表" }
func (s *ProducerList) Help() []base.CommandOption { return []base.CommandOption{} }

func (s *ProducerList) Execute(_ map[string]string) {
	producers := facade.Queue().Producers()
	if len(producers) == 0 {
		color.Yellow("no registered producers")
		return
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
}

func init() { cli.Register(&ProducerList{}) }
