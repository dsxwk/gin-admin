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

type MakeService struct {
	base.BaseCommand
}

func (m *MakeService) Name() string {
	return "make:service"
}

func (m *MakeService) Description() string {
	return "服务创建"
}

func (m *MakeService) Help() []base.CommandOption {
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
				Short:   "F",
				Long:    "function",
				Default: "FuncName",
			},
			"方法名称, 如: list",
			false,
		},
		{
			base.Flag{
				Short:   "d",
				Long:    "desc",
				Default: "service-desc",
			},
			"描述, 如: 列表",
			false,
		},
	}
}

func (m *MakeService) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())
	color.Green("执行命令: %s %s", m.Name(), m.FormatArgs(values))
	_make := strings.TrimPrefix(m.Name(), "make:")
	f := m.GetMakeFile(values["file"], _make)
	m.generateFile(_make, f, values["function"], values["desc"])
}

func init() {
	cli.Register(&MakeService{})
}

func (m *MakeService) generateFile(_make, file, function, desc string) {
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
		Function    string // 如果为空,使用默认值
		Description string // 如果为空,使用默认值
	}{
		Package:     packageName,
		Name:        pkg.ToUpperCamel(strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))),
		Function:    pkg.ToUpperCamel(function),
		Description: desc,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		color.Red("Error executing template:", err.Error())
		os.Exit(1)
	}

	color.Green(flag.Success + " 服务文件: " + file + " 生成成功!")
}
