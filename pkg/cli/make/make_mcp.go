package make

import (
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

// MakeMcp MCP工具创建命令
type MakeMcp struct {
	base.BaseCommand
}

func (m *MakeMcp) Name() string {
	return "make:mcp"
}

func (m *MakeMcp) Description() string {
	return "MCP工具创建"
}

func (m *MakeMcp) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "f",
				Long:  "file",
			},
			"文件路径, 如: article_search",
			true,
		},
		{
			base.Flag{
				Short: "n",
				Long:  "name",
			},
			"工具名称, 如: search_articles, 默认取文件名的蛇形命名",
			false,
		},
		{
			base.Flag{
				Short:   "d",
				Long:    "desc",
				Default: "工具描述",
			},
			"工具描述, 如: 根据关键词搜索文章",
			false,
		},
	}
}

func (m *MakeMcp) Execute(values map[string]string) {
	_make := strings.TrimPrefix(m.Name(), "make:")
	file := m.GetMakeFile(values["file"], _make)
	m.generateFile(_make, file, values["file"], values["name"], values["desc"])
}

func init() {
	cli.Register(&MakeMcp{})
}

func (m *MakeMcp) generateFile(_make, file, rawName, name, desc string) {
	templateFile := m.GetTemplate(_make)
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		flag.Errorf("Error parsing template: %s", err.Error())
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
		flag.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("MCP工具文件: %s 生成成功!", file)
}
