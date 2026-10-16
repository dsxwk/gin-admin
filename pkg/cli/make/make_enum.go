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

type MakeEnum struct {
	base.Command
}

// NewMakeEnum 创建命令
func NewMakeEnum() *MakeEnum {
	return &MakeEnum{
		Command: base.NewCommand(base.Config{
			Name:        "make:enum",
			Description: "枚举创建",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "f", Long: "file"},
					Desc:     "文件路径, 如: test/user",
					Required: true,
				},
				{
					Flag:     base.Flag{Short: "d", Long: "desc", Default: "enum-desc"},
					Desc:     "描述, 如: enum-desc",
					Required: false,
				},
			},
		}),
	}
}

func (m *MakeEnum) Handle(values map[string]string) error {
	_make := strings.TrimPrefix(m.Name(), "make:")
	f := m.GetMakeFile(values["file"], _make)
	m.generateFile(_make, f, values["desc"])
	return nil
}

func (m *MakeEnum) generateFile(_make, file, desc string) {
	templateFile := m.GetTemplate(_make)
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		terminal.Errorf("Error parsing template: %s", err.Error())
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
		Description: desc,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		terminal.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	terminal.Successf("命令行文件: " + file + " 生成成功!")
}
