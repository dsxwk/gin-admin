package make

import (
	"fmt"
	"gin/pkg/cli"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type MakeListener struct {
	base.Command
}

// NewMakeListener 创建命令
func NewMakeListener() *MakeListener {
	return &MakeListener{
		Command: base.NewCommand(base.Config{
			Name:        "make:listener",
			Description: "创建监听",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "f", Long: "file"},
					Desc:     "文件路径, 如: login/test",
					Required: true,
				},
				{
					Flag:     base.Flag{Short: "e", Long: "event"},
					Desc:     "事件数据, 如: UserLogin",
					Required: true,
				},
			},
		}),
	}
}

func (m *MakeListener) Handle(values map[string]string) error {
	_make := "listener"
	file := m.GetMakeFile(values["file"], _make)
	eventName := values["event"]
	m.generateFile(_make, file, eventName)
	return nil
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
