package mcp

import (
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg/cli"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

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

func (s *McpList) Execute(_ map[string]string) {
	tools := facade.MCP().Tools()
	if len(tools) == 0 {
		color.Yellow("暂无注册的MCP工具")
		return
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
}

func init() {
	cli.Register(&McpList{})
}
