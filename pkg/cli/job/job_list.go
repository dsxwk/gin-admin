package job

import (
	"fmt"
	"gin/common/base"
	"gin/pkg/cli"
	"gin/pkg/serviceprovider/job"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"sort"
	"strings"
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
	jobs := job.GetAll()
	if len(jobs) == 0 {
		color.Yellow("暂无注册的Job")
		return
	}

	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].Name() < jobs[j].Name()
	})

	maxNameLen := 0
	maxConnLen := 0
	maxDescLen := 0
	maxRetryLen := 0
	maxDelayLen := 0
	for _, jb := range jobs {
		nameLen := runewidth.StringWidth(jb.Name())
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}

		conn := jb.Connection()
		if conn == "" {
			conn = "redis"
		}
		connLen := runewidth.StringWidth(conn)
		if connLen > maxConnLen {
			maxConnLen = connLen
		}

		descLen := runewidth.StringWidth(jb.Description())
		if descLen > maxDescLen {
			maxDescLen = descLen
		}

		retryStr := fmt.Sprintf("%d", jb.Retry())
		retryLen := runewidth.StringWidth(retryStr)
		if retryLen > maxRetryLen {
			maxRetryLen = retryLen
		}

		delayStr := fmt.Sprintf("%dms", jb.Delay())
		delayLen := runewidth.StringWidth(delayStr)
		if delayLen > maxDelayLen {
			maxDelayLen = delayLen
		}
	}

	titleNameLen := runewidth.StringWidth("Job名称")
	titleConnLen := runewidth.StringWidth("连接")
	titleRetryLen := runewidth.StringWidth("重试次数")
	titleDelayLen := runewidth.StringWidth("延迟(ms)")
	titleDescLen := runewidth.StringWidth("描述")
	if titleNameLen > maxNameLen {
		maxNameLen = titleNameLen
	}
	if titleConnLen > maxConnLen {
		maxConnLen = titleConnLen
	}
	if titleRetryLen > maxRetryLen {
		maxRetryLen = titleRetryLen
	}
	if titleDelayLen > maxDelayLen {
		maxDelayLen = titleDelayLen
	}
	if titleDescLen > maxDescLen {
		maxDescLen = titleDescLen
	}

	totalWidth := maxNameLen + maxConnLen + maxRetryLen + maxDelayLen + maxDescLen + 16

	color.Yellow("┌" + strings.Repeat("─", totalWidth-2) + "┐")

	titleLine := fmt.Sprintf("│ %s   %s   %s   %s   %s "+color.YellowString("│"),
		color.HiWhiteString(padRight("Job名称", maxNameLen)),
		color.HiWhiteString(padRight("连接", maxConnLen)),
		color.HiWhiteString(padRight("重试次数", maxRetryLen)),
		color.HiWhiteString(padRight("延迟(ms)", maxDelayLen)),
		color.HiWhiteString(padRight("描述", maxDescLen)))
	color.Yellow(titleLine)

	color.Yellow("├" + strings.Repeat("─", totalWidth-2) + "┤")

	for _, jb := range jobs {
		conn := jb.Connection()
		if conn == "" {
			conn = "redis"
		}

		retryStr := fmt.Sprintf("%d", jb.Retry())
		delayStr := fmt.Sprintf("%dms", jb.Delay())

		contentLine := fmt.Sprintf("│ %s   %s   %s   %s   %s "+color.YellowString("│"),
			color.GreenString(padRight(jb.Name(), maxNameLen)),
			color.CyanString(padRight(conn, maxConnLen)),
			color.YellowString(padRight(retryStr, maxRetryLen)),
			color.MagentaString(padRight(delayStr, maxDelayLen)),
			color.WhiteString(padRight(jb.Description(), maxDescLen)))
		color.Yellow(contentLine)
	}

	color.Yellow("└" + strings.Repeat("─", totalWidth-2) + "┘")

	color.Cyan(fmt.Sprintf("总计 %d 个Job\n", len(jobs)))
}

func padRight(s string, width int) string {
	currentWidth := runewidth.StringWidth(s)
	if currentWidth >= width {
		return s
	}
	padding := width - currentWidth
	return s + strings.Repeat(" ", padding)
}

func init() {
	cli.Register(&JobList{})
}
