package make

import (
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/cli"
	"github.com/fatih/color"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type MakeMiddleware struct {
	base.BaseCommand
}

func (m *MakeMiddleware) Name() string {
	return "make:middleware"
}

func (m *MakeMiddleware) Description() string {
	return "中间件创建"
}

func (m *MakeMiddleware) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "f",
				Long:  "file",
			},
			"文件路径, 如: v1/user",
			true,
		},
		{
			base.Flag{
				Short:   "d",
				Long:    "desc",
				Default: "middleware-desc",
			},
			"描述, 如: 列表",
			false,
		},
	}
}

func (m *MakeMiddleware) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())
	color.Green("执行命令: %s %s", m.Name(), m.FormatArgs(values))
	_make := strings.TrimPrefix(m.Name(), "make:")
	f := m.GetMakeFile(values["file"], _make)
	m.generateFile(_make, f, values["desc"])
}

func init() {
	cli.Register(&MakeMiddleware{})
}

func (m *MakeMiddleware) generateFile(_make, file, desc string) {
	templateFile := m.GetTemplate(_make)
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		color.Red("Error parsing template:", err.Error())
		os.Exit(1)
	}

	// 提取包名 (文件路径中的最后一个目录作为包名)
	packageName := filepath.Base(filepath.Dir(file))

	// 创建文件
	f := m.CheckDirAndFile(file)
	if f == nil {
		return
	}

	data := struct {
		Package     string // 提取的包名
		Name        string // 模块名称(首字母大写)
		Description string // 如果为空,使用默认值
	}{
		Package:     packageName,
		Name:        pkg.ToUpperCamel(strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))),
		Description: desc,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		color.Red("Error executing template:", err.Error())
		os.Exit(1)
	}

	color.Green(flag.Success + "  中间件文件: " + file + " 生成成功!")
}
