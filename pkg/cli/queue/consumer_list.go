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

type ConsumerList struct{}

func (s *ConsumerList) Name() string {
	return "consumer:list"
}

func (s *ConsumerList) Description() string {
	return "消费者列表"
}

func (s *ConsumerList) Help() []base.CommandOption {
	return []base.CommandOption{}
}

func (s *ConsumerList) Execute(args []string) {
	consumers := queue.GetConsumerRegistry().GetAll()
	if len(consumers) == 0 {
		color.Yellow("暂无注册的消费者")
		return
	}

	// 按名称排序
	sort.Slice(consumers, func(i, j int) bool {
		return consumers[i].Name() < consumers[j].Name()
	})

	// 计算最大名称宽度和描述宽度(使用显示宽度)
	maxNameLen := 0
	maxDescLen := 0
	for _, c := range consumers {
		nameLen := runewidth.StringWidth(c.Name())
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}

		descLen := runewidth.StringWidth(c.Description())
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
	titleNameLen := runewidth.StringWidth("消费者名称")
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
		color.HiWhiteString(padRight("消费者名称", maxNameLen)),
		color.HiWhiteString(padRight("描述", maxDescLen)))
	color.Yellow(titleLine)

	// 打印分隔线
	color.Yellow("├" + strings.Repeat("─", maxNameLen+2) + "─" + strings.Repeat("─", maxDescLen+2) + "┤")

	// 打印消费者列表
	for _, c := range consumers {
		coloredName := color.GreenString(padRight(c.Name(), maxNameLen))
		coloredDesc := color.WhiteString(padRight(c.Description(), maxDescLen))
		contentLine := fmt.Sprintf("│ %s   %s "+color.YellowString("│"), coloredName, coloredDesc)
		color.Yellow(contentLine)
	}

	// 打印底部边框
	color.Yellow("└" + strings.Repeat("─", totalWidth-2) + "┘")

	// 打印统计信息
	color.Cyan(fmt.Sprintf("总计 %d 个消费者\n", len(consumers)))
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
	cli.Register(&ConsumerList{})
}
