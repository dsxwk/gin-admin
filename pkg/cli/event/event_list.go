package event

import (
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg/cli"
	"gin/pkg/eventbus"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"sort"
	"strings"
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

func (s *EventList) Execute(args []string) {
	list := facade.Event[eventbus.Event]().List()
	if len(list) == 0 {
		color.Yellow("暂无注册的事件")
		return
	}

	// 按名称排序
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	// 计算最大名称宽度和描述宽度(使用显示宽度)
	maxNameLen := 0
	maxDescLen := 0
	for _, v := range list {
		nameLen := runewidth.StringWidth(v.Name)
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}

		descLen := runewidth.StringWidth(v.Description)
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
	titleNameLen := runewidth.StringWidth("事件名称")
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
		color.HiWhiteString(padRight("事件名称", maxNameLen)),
		color.HiWhiteString(padRight("描述", maxDescLen)))
	color.Yellow(titleLine)

	// 打印分隔线
	color.Yellow("├" + strings.Repeat("─", maxNameLen+2) + "─" + strings.Repeat("─", maxDescLen+2) + "┤")

	// 打印事件列表
	for _, v := range list {
		coloredName := color.GreenString(padRight(v.Name, maxNameLen))
		coloredDesc := color.WhiteString(padRight(v.Description, maxDescLen))
		contentLine := fmt.Sprintf("│ %s   %s "+color.YellowString("│"), coloredName, coloredDesc)
		color.Yellow(contentLine)
	}

	// 打印底部边框
	color.Yellow("└" + strings.Repeat("─", totalWidth-2) + "┘")

	// 打印统计信息
	color.Cyan(fmt.Sprintf("总计 %d 个事件\n", len(list)))
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
	cli.Register(&EventList{})
}
