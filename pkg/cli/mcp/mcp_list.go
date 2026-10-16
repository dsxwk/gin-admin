package mcp

import (
	"fmt"
	"gin/app/facade"
	"gin/pkg/cli"
	"gin/pkg/cli/base"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type McpList struct{ base.Command }

// NewMcpList 创建命令
func NewMcpList() *McpList {
	return &McpList{
		Command: base.NewCommand(base.Config{
			Name:        "mcp:list",
			Description: "MCP工具列表",
		}),
	}
}

func (s *McpList) Handle(_ map[string]string) error {
	tools := facade.MCP().Tools()
	if len(tools) == 0 {
		color.Yellow("暂无注册的MCP工具")
		return nil

	}

	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Name() < tools[j].Name()
	})

	writer := cli.NewTable()
	writer.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Colors: text.Colors{text.FgGreen}},
		{Number: 2, Colors: text.Colors{text.FgWhite}, WidthMax: 60},
	})
	writer.AppendHeader(table.Row{"工具名称", "描述"})

	for _, tool := range tools {
		writer.AppendRow(table.Row{
			tool.Name(),
			strings.TrimSpace(tool.Description()),
		})
	}

	fmt.Println(writer.Render())
	color.Cyan("总计 %d 个MCP工具\n", len(tools))
	return nil
}
