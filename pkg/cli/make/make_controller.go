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

type MakeController struct {
	base.BaseCommand
}

func (m *MakeController) Name() string {
	return "make:controller"
}

func (m *MakeController) Description() string {
	return "控制器创建"
}

func (m *MakeController) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "f",
				Long:  "file",
			},
			"文件路径, 如: v1/user",
			true,
		},
		{
			base.Flag{
				Short:   "F",
				Long:    "function",
				Default: "FuncName",
			},
			"方法名称, 如: list",
			false,
		},
		{
			base.Flag{
				Short:   "m",
				Long:    "method",
				Default: "get",
			},
			"请求方式, 如: get",
			false,
		},
		{
			base.Flag{
				Short:   "r",
				Long:    "router",
				Default: "/your/router",
			},
			"路由地址, 如: /user",
			false,
		},
		{
			base.Flag{
				Short:   "d",
				Long:    "desc",
				Default: "func-desc",
			},
			"描述, 如: 列表",
			false,
		},
	}
}

func (m *MakeController) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())
	_make := strings.TrimPrefix(m.Name(), "make:")
	f := m.GetMakeFile(values["file"], _make)
	m.generateFile(_make, f, values["function"], values["method"], values["router"], values["desc"])
}

func init() {
	cli.Register(&MakeController{})
}

func (m *MakeController) generateFile(_make, file, function, method, router, desc string) {
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
		Function    string // 如果为空,使用默认值
		Router      string // 如果为空,使用默认值
		Method      string // 如果为空,使用默认值
		Description string // 如果为空,使用默认值
	}{
		Package:     packageName,
		Name:        pkg.ToUpperCamel(strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))),
		Function:    pkg.UcFirst(function),
		Router:      router,
		Method:      method,
		Description: desc,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		flag.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("控制器文件: " + file + " 生成成功!")
}
