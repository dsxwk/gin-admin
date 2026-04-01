package make

import (
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/cli"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

type MakeRouter struct {
	base.BaseCommand
}

func (m *MakeRouter) Name() string {
	return "make:router"
}

func (m *MakeRouter) Description() string {
	return "路由创建"
}

func (m *MakeRouter) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "f",
				Long:  "file",
			},
			"文件路径, 如: user",
			true,
		},
		{
			base.Flag{
				Short:   "d",
				Long:    "desc",
				Default: "router-desc",
			},
			"路由描述, 如: 用户路由",
			false,
		},
	}
}

func (m *MakeRouter) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())
	_make := strings.TrimPrefix(m.Name(), "make:")
	f := m.GetMakeFile(values["file"], _make)
	m.generateFile(_make, f, values["desc"])
}

func init() {
	cli.Register(&MakeRouter{})
}

func (m *MakeRouter) generateFile(_make, file, desc string) {
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
		Description string // 如果为空,使用默认值
	}{
		Package:     packageName,
		Name:        pkg.ToUpperCamel(strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))),
		Description: desc,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		flag.Successf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("路由文件: " + file + " 生成成功!")
}
