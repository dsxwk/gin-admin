package make

import (
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/cli"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type MakeFacade struct {
	base.BaseCommand
}

func (m *MakeFacade) Name() string {
	return "make:facade"
}

func (m *MakeFacade) Description() string {
	return "创建门面"
}

func (m *MakeFacade) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "f",
				Long:  "file",
			},
			"文件路径, 如: cache",
			true,
		},
		{
			base.Flag{
				Short:   "d",
				Long:    "desc",
				Default: "facade-desc",
			},
			"门面描述, 如: 缓存门面",
			false,
		},
	}
}

func (m *MakeFacade) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())

	_make := strings.TrimPrefix(m.Name(), "make:")
	fileName := values["file"]

	file := m.GetMakeFile(fileName, _make)
	m.generateFacade(_make, file, values["desc"])
}

func init() {
	cli.Register(&MakeFacade{})
}

func (m *MakeFacade) generateFacade(_make, file, facadeDesc string) {
	templateFile := m.GetTemplate(_make)
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		flag.Errorf("Error parsing template: %s", err.Error())
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
	facadeName := pkg.ToUpperCamel(strings.ReplaceAll(baseName, "_", " "))
	facadeName = strings.ReplaceAll(facadeName, " ", "")

	// 生成小驼峰变量名
	facadeVar := pkg.SnakeToCamel(facadeName)

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
		flag.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("门面文件: " + file + " 生成成功!")
}
