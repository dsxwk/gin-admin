package make

import (
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/samber/lo"
)

type MakeFacade struct {
	base.Command
}

// NewMakeFacade 创建命令
func NewMakeFacade() *MakeFacade {
	return &MakeFacade{
		Command: base.NewCommand(base.Config{
			Name:        "make:facade",
			Description: "创建门面",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "f", Long: "file"},
					Desc:     "文件路径, 如: cache",
					Required: true,
				},
				{
					Flag:     base.Flag{Short: "d", Long: "desc", Default: "facade-desc"},
					Desc:     "门面描述, 如: 缓存门面",
					Required: false,
				},
			},
		}),
	}
}

func (m *MakeFacade) Handle(values map[string]string) error {

	_make := strings.TrimPrefix(m.Name(), "make:")
	fileName := values["file"]

	file := m.GetMakeFile(fileName, _make)
	m.generateFacade(_make, file, values["desc"])
	return nil
}

func (m *MakeFacade) generateFacade(_make, file, facadeDesc string) {
	templateFile := m.GetTemplate(_make)
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		terminal.Errorf("Error parsing template: %s", err.Error())
		os.Exit(1)
	}

	// 创建文件
	f := m.CheckDirAndFile(file)
	if f == nil {
		return
	}

	// 从文件名提取门面名称
	baseName := strings.TrimSuffix(filepath.Base(file), ".go")

	// 转换为大驼峰作为门面名称
	facadeName := lo.PascalCase(baseName)

	// 生成小驼峰变量名
	facadeVar := lo.CamelCase(facadeName)

	data := struct {
		Package    string // 包名(facade)
		FacadeName string // 门面名称(大驼峰)
		FacadeVar  string // 门面变量名(小写驼峰)
		Desc       string // 门面描述
	}{
		Package:    "facade",
		FacadeName: facadeName,
		FacadeVar:  facadeVar,
		Desc:       facadeDesc,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		terminal.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	terminal.Successf("门面文件: " + file + " 生成成功!")
}
