package mcp

import (
	"fmt"
	"gin/common/base"
	"gin/pkg/cli"
	svcmcp "gin/pkg/serviceprovider/mcp"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"sort"
	"strings"
)

// toolRow MCP工具行数据
type toolRow struct {
	name      string
	descLines []string
}

type McpList struct{}

func (s *McpList) Name() string {
	return "mcp:list"
}

func (s *McpList) Description() string {
	return "MCP工具列表"
}

func (s *McpList) Help() []base.CommandOption {
	return []base.CommandOption{}
}

func (s *McpList) Execute(values map[string]string) {
	tools := svcmcp.GetAll()
	if len(tools) == 0 {
		color.Yellow("暂无注册的MCP工具")
		return
	}

	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Name() < tools[j].Name()
	})

	maxNameLen := runewidth.StringWidth("工具名称")
	maxDescLen := runewidth.StringWidth("描述")

	rows := make([]toolRow, 0, len(tools))
	for _, t := range tools {
		descLines := s.splitDesc(t.Description())
		rows = append(rows, toolRow{name: t.Name(), descLines: descLines})

		nameLen := runewidth.StringWidth(t.Name())
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}

		for _, line := range descLines {
			lineLen := runewidth.StringWidth(line)
			if lineLen > maxDescLen {
				maxDescLen = lineLen
			}
		}
	}

	totalWidth := maxNameLen + maxDescLen + 6

	color.Yellow("┌" + strings.Repeat("─", totalWidth-2) + "┐")

	titleLine := fmt.Sprintf("│%s   %s "+color.YellowString("│"),
		color.HiWhiteString(padRight("工具名称", maxNameLen)),
		color.HiWhiteString(padRight("描述", maxDescLen)))
	color.Yellow(titleLine)

	color.Yellow("├" + strings.Repeat("─", totalWidth-2) + "┤")

	for _, row := range rows {
		for i, line := range row.descLines {
			name := row.name
			if i > 0 {
				name = ""
			}

			contentLine := fmt.Sprintf("│%s   %s "+color.YellowString("│"),
				color.GreenString(padRight(name, maxNameLen)),
				color.WhiteString(padRight(line, maxDescLen)))
			color.Yellow(contentLine)
		}
	}

	color.Yellow("└" + strings.Repeat("─", totalWidth-2) + "┘")

	color.Cyan(fmt.Sprintf("总计 %d 个MCP工具\n", len(tools)))
}

// splitDesc 将描述按换行拆分并去除首尾空白
func (s *McpList) splitDesc(desc string) []string {
	rawLines := strings.Split(desc, "\n")
	lines := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) == 0 {
		lines = append(lines, "")
	}
	return lines
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
	cli.Register(&McpList{})
}
