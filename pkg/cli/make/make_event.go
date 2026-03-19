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

type MakeEvent struct {
	base.BaseCommand
}

func (m *MakeEvent) Name() string {
	return "make:event"
}

func (m *MakeEvent) Description() string {
	return "创建事件"
}

func (m *MakeEvent) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "f",
				Long:  "file",
			},
			"文件路径, 如: login/test",
			true,
		},
		{
			base.Flag{
				Short: "n",
				Long:  "name",
			},
			"事件名称, 如: test-event",
			false,
		},
		{
			base.Flag{
				Short: "d",
				Long:  "desc",
			},
			"事件描述, 如: 测试事件",
			false,
		},
	}
}

func (m *MakeEvent) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())
	color.Green("执行命令: %s %v", m.Name(), values)
	_make := "event"
	file := m.GetMakeFile(values["file"], _make)
	name := values["name"]
	desc := values["desc"]
	m.generateFile(_make, file, name, desc)
}

func init() {
	cli.Register(&MakeEvent{})
}

func (m *MakeEvent) generateFile(_make, file, name, desc string) {
	templateFile := m.GetTemplate(_make)
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		color.Red("Error parsing template:", err.Error())
		os.Exit(1)
	}

	packageName := filepath.Base(filepath.Dir(file))
	f := m.CheckDirAndFile(file)
	if f == nil {
		return
	}

	data := struct {
		Package     string
		Struct      string
		Name        string
		Description string
	}{
		Package:     packageName,
		Struct:      pkg.ToUpperCamel(strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))),
		Name:        name,
		Description: desc,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		color.Red("Error executing template:", err.Error())
		os.Exit(1)
	}

	color.Green(flag.Success+"  事件文件: %s 生成成功!", file)
}
