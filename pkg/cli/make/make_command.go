package make

import (
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"github.com/samber/lo"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type MakeCommand struct {
	base.BaseCommand
}

func (m *MakeCommand) Name() string {
	return "make:command"
}

func (m *MakeCommand) Description() string {
	return "命令创建"
}

func (m *MakeCommand) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "f",
				Long:  "file",
			},
			"文件路径, 如: cronjob/demo",
			true,
		},
		{
			base.Flag{
				Short: "n",
				Long:  "name",
			},
			"命令名称, 如: demo-test",
			false,
		},
		{
			base.Flag{
				Short:   "d",
				Long:    "desc",
				Default: "command-desc",
			},
			"描述, 如: command-desc",
			false,
		},
	}
}

func (m *MakeCommand) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())
	_make := strings.TrimPrefix(m.Name(), "make:")
	f := m.GetMakeFile(values["file"], _make)
	m.generateFile(_make, f, values["name"], values["desc"])
}

func init() {
	cli.Register(&MakeCommand{})
}

func (m *MakeCommand) generateFile(_make, file, name, desc string) {
	templateFile := m.GetTemplate(_make)
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		flag.Errorf("Error parsing template: %s", err.Error())
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
		Command     string // 如果为空,使用默认值
		Description string // 如果为空,使用默认值
	}{
		Package:     packageName,
		Name:        lo.PascalCase(strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))),
		Command:     name,
		Description: desc,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		flag.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("命令行文件: " + file + " 生成成功!")
}
