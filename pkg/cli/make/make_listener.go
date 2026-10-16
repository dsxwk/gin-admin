package make

import (
	"fmt"
	"gin/common/base"
	"gin/pkg/terminal"
	"gin/pkg/cli"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
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

func (m *MakeListener) Execute(values map[string]string) {
	_make := "listener"
	file := m.GetMakeFile(values["file"], _make)
	eventName := values["event"]
	m.generateFile(_make, file, eventName)
}

func (m *MakeListener) generateFile(_make, file, eventName string) {
	templateFile := m.GetTemplate(_make)
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		terminal.Errorf("Error parsing template: %s", err.Error())
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
		Name:      lo.PascalCase(strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))),
		EventName: eventName,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		terminal.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	listenersFile := filepath.Join("app", "listener", "registry.go")
	qualifier := ""
	listenerDir := filepath.ToSlash(filepath.Clean(filepath.Dir(file)))
	if listenerDir != "app/listener" {
		alias, importErr := cli.AddImport(listenersFile, cli.ImportPath(file))
		if importErr != nil {
			terminal.Errorf("自动添加监听器导入失败: %s", importErr.Error())
			return
		}
		qualifier = alias + "."
	}

	item := fmt.Sprintf(
		"eventbus.Bind(&%s%sListener{}, event.%s{})",
		qualifier,
		data.Name,
		eventName,
	)
	if err = cli.AddItem(listenersFile, "return []eventbus.ListenerRegistration{", item); err != nil {
		terminal.Errorf("自动注册监听器失败: %s", err.Error())
		return
	}

	terminal.Successf("监听文件: %s 生成成功!", file)
}
