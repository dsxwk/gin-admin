package make

import (
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type MakeEvent struct {
	base.Command
}

// NewMakeEvent 创建命令
func NewMakeEvent() *MakeEvent {
	return &MakeEvent{
		Command: base.NewCommand(base.Config{
			Name:        "make:event",
			Description: "创建事件",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "f", Long: "file"},
					Desc:     "文件路径, 如: login/test",
					Required: true,
				},
				{
					Flag:     base.Flag{Short: "n", Long: "name"},
					Desc:     "事件名称, 如: test-event",
					Required: false,
				},
				{
					Flag:     base.Flag{Short: "d", Long: "desc"},
					Desc:     "事件描述, 如: 测试事件",
					Required: false,
				},
			},
		}),
	}
}

func (m *MakeEvent) Handle(values map[string]string) error {
	_make := "event"
	file := m.GetMakeFile(values["file"], _make)
	name := values["name"]
	desc := values["desc"]
	m.generateFile(_make, file, name, desc)
	return nil
}

func (m *MakeEvent) generateFile(_make, file, name, desc string) {
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
		Package     string
		Struct      string
		Name        string
		Description string
	}{
		Package:     packageName,
		Struct:      lo.PascalCase(strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))),
		Name:        name,
		Description: desc,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		terminal.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	terminal.Successf("事件文件: %s 生成成功!", file)
}
