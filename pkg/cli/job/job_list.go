package job

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

type JobList struct{}

func (s *JobList) Name() string {
	return "job:list"
}

func (s *JobList) Description() string {
	return "Job列表"
}

func (s *JobList) Help() []base.CommandOption {
	return []base.CommandOption{}
}

func (s *JobList) Execute(values map[string]string) {
	jobs := facade.Job().List()
	if len(jobs) == 0 {
		color.Yellow("暂无注册的Job")
		return
	}

	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Name() < jobs[j].Name()
	})

	writer := cli.NewTable()
	writer.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Colors: text.Colors{text.FgGreen}},
		{Number: 2, Colors: text.Colors{text.FgCyan}},
		{Number: 3, Colors: text.Colors{text.FgYellow}},
		{Number: 4, Colors: text.Colors{text.FgMagenta}},
		{Number: 5, Colors: text.Colors{text.FgWhite}},
	})
	writer.AppendHeader(table.Row{"Job名称", "连接", "重试次数", "延迟(ms)", "描述"})

	for _, jb := range jobs {
		conn := jb.Connection()
		if conn == "" {
			conn = "redis"
		}

		retryStr := fmt.Sprintf("%d", jb.Retry())
		delayStr := fmt.Sprintf("%dms", jb.Delay())

		writer.AppendRow(table.Row{
			jb.Name(),
			conn,
			retryStr,
			delayStr,
			jb.Description(),
		})
	}

	fmt.Println(writer.Render())
	color.Cyan(fmt.Sprintf("总计 %d 个Job\n", len(jobs)))
}

func init() {
	cli.Register(&JobList{})
}
