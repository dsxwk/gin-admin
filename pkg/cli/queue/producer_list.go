package queue

import (
	"fmt"
	"gin/common/base"
	"gin/pkg/cli"
	"gin/pkg/queue"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"sort"
	"strings"
)

type ProducerList struct{}

func (s *ProducerList) Name() string {
	return "producer:list"
}

func (s *ProducerList) Description() string {
	return "生产者列表"
}

func (s *ProducerList) Help() []base.CommandOption {
	return []base.CommandOption{}
}

func (s *ProducerList) Execute(args []string) {
	producers := queue.GetProducerRegistry().GetAll()
	if len(producers) == 0 {
		color.Yellow("暂无注册的生产者")
		return
	}

	// 按名称排序
	sort.Slice(producers, func(i, j int) bool {
		return producers[i].Name() < producers[j].Name()
	})

	// 计算最大名称宽度和描述宽度(使用显示宽度)
	maxNameLen := 0
	maxDescLen := 0
	for _, p := range producers {
		nameLen := runewidth.StringWidth(p.Name())
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}

		descLen := runewidth.StringWidth(p.Description())
		if descLen > maxDescLen {
			maxDescLen = descLen
		}
	}

	// 设置最小宽度
	if maxNameLen < 20 {
		maxNameLen = 20
	}
	if maxDescLen < 35 {
		maxDescLen = 35
	}

	// 计算标题的显示宽度
	titleNameLen := runewidth.StringWidth("生产者名称")
	titleDescLen := runewidth.StringWidth("描述")
	if titleNameLen > maxNameLen {
		maxNameLen = titleNameLen
	}
	if titleDescLen > maxDescLen {
		maxDescLen = titleDescLen
	}

	// 计算总宽度(名称宽度+描述宽度+边框)
	totalWidth := maxNameLen + maxDescLen + 7

	// 打印顶部边框
	color.Yellow("┌" + strings.Repeat("─", totalWidth-2) + "┐")

	// 打印标题行
	titleLine := fmt.Sprintf("│ %s   %s "+color.YellowString("│"),
		color.HiWhiteString(padRight("生产者名称", maxNameLen)),
		color.HiWhiteString(padRight("描述", maxDescLen)))
	color.Yellow(titleLine)

	// 打印分隔线
	color.Yellow("├" + strings.Repeat("─", maxNameLen+2) + "─" + strings.Repeat("─", maxDescLen+2) + "┤")

	// 打印生产者列表
	for _, c := range producers {
		coloredName := color.GreenString(padRight(c.Name(), maxNameLen))
		coloredDesc := color.WhiteString(padRight(c.Description(), maxDescLen))
		contentLine := fmt.Sprintf("│ %s   %s "+color.YellowString("│"), coloredName, coloredDesc)
		color.Yellow(contentLine)
	}

	// 打印底部边框
	color.Yellow("└" + strings.Repeat("─", totalWidth-2) + "┘")

	// 打印统计信息
	color.Cyan(fmt.Sprintf("总计 %d 个生产者\n", len(producers)))
}

func init() {
	cli.Register(&ProducerList{})
}
