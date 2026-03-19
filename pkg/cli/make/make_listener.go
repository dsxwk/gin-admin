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

type MakeListener struct {
	base.BaseCommand
}

func (m *MakeListener) Name() string {
	return "make:listener"
}

func (m *MakeListener) Description() string {
	return "创建监听"
}

func (m *MakeListener) Help() []base.CommandOption {
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
				Short: "e",
				Long:  "event",
			},
			"事件数据, 如: UserLogin",
			true,
		},
	}
}

func (m *MakeListener) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())
	color.Green("执行命令: %s %v", m.Name(), values)
	_make := "listener"
	file := m.GetMakeFile(values["file"], _make)
	eventName := values["event"]
	m.generateFile(_make, file, eventName)
}

func init() {
	cli.Register(&MakeListener{})
}

func (m *MakeListener) generateFile(_make, file, eventName string) {
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
		Package   string
		Name      string
		EventName string
	}{
		Package:   packageName,
		Name:      pkg.ToUpperCamel(strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))),
		EventName: eventName,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		color.Red("Error executing template:", err.Error())
		os.Exit(1)
	}

	color.Green(flag.Success+" 监听文件: %s 生成成功!", file)
}
