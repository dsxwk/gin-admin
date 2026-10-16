package make

import (
	"gin/pkg/cli"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

// MakeMcp MCP工具创建命令
type MakeMcp struct {
	base.Command
}

// NewMakeMcp 创建命令
func NewMakeMcp() *MakeMcp {
	return &MakeMcp{
		Command: base.NewCommand(base.Config{
			Name:        "make:mcp",
			Description: "MCP工具创建",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "f", Long: "file"},
					Desc:     "文件路径, 如: article_search",
					Required: true,
				},
				{
					Flag:     base.Flag{Short: "n", Long: "name"},
					Desc:     "工具名称, 如: search_articles, 默认取文件名的蛇形命名",
					Required: false,
				},
				{
					Flag:     base.Flag{Short: "d", Long: "desc", Default: "工具描述"},
					Desc:     "工具描述, 如: 根据关键词搜索文章",
					Required: false,
				},
			},
		}),
	}
}

func (m *MakeMcp) Handle(values map[string]string) error {
	_make := strings.TrimPrefix(m.Name(), "make:")
	file := m.GetMakeFile(values["file"], _make)
	m.generateFile(_make, file, values["file"], values["name"], values["desc"])
	return nil
}

func (m *MakeMcp) generateFile(_make, file, rawName, name, desc string) {
	templateFile := m.GetTemplate(_make)
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		terminal.Errorf("Error parsing template: %s", err.Error())
		os.Exit(1)
	}

	// 文件名(去掉路径和.go后缀)
	baseName := path.Base(strings.TrimSuffix(filepath.ToSlash(rawName), ".go"))

	// 结构体名
	structName := lo.PascalCase(baseName)

	// 工具名(默认蛇形命名)
	toolName := name
	if toolName == "" {
		toolName = lo.SnakeCase(baseName)
	}

	data := struct {
		Package     string
		Name        string
		ToolName    string
		Description string
	}{
		Package:     "mcp",
		Name:        structName,
		ToolName:    toolName,
		Description: desc,
	}

	f := m.CheckDirAndFile(file)
	if f == nil {
		return
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		terminal.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	toolsFile := filepath.Join("app", "mcp", "registry.go")
	qualifier := ""
	mcpDir := filepath.ToSlash(filepath.Clean(filepath.Dir(file)))
	if mcpDir != "app/mcp" {
		alias, importErr := cli.AddImport(toolsFile, cli.ImportPath(file))
		if importErr != nil {
			terminal.Errorf("自动添加MCP工具导入失败: %s", importErr.Error())
			return
		}
		qualifier = alias + "."
	}
	if err = cli.AddItem(toolsFile, "return []servicemcp.Tool{", qualifier+"New"+structName+"Tool()"); err != nil {
		terminal.Errorf("自动注册MCP工具失败: %s", err.Error())
		return
	}

	terminal.Successf("MCP工具文件: %s 生成成功!", file)
}
