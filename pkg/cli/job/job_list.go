package job

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

type JobList struct{ base.Command }

// NewJobList 创建命令
func NewJobList() *JobList {
	return &JobList{
		Command: base.NewCommand(base.Config{
			Name:        "job:list",
			Description: "Job列表",
		}),
	}
}

func (s *JobList) Handle(_ map[string]string) error {
	jobs := facade.Job().List()
	if len(jobs) == 0 {
		color.Yellow("暂无注册的Job")
		return nil

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
	return nil
}
