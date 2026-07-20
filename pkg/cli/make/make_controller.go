package make

import (
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"github.com/samber/lo"
	"html/template"
	"os"
	"path"
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
				Short:   "d",
				Long:    "desc",
				Default: "用户",
			},
			"描述, 如: 用户",
			false,
		},
	}
}

func (m *MakeController) Execute(values map[string]string) {
	_make := strings.TrimPrefix(m.Name(), "make:")
	f := m.GetMakeFile(values["file"], _make)
	m.generateFile(_make, f, values["desc"])
}

func init() {
	cli.Register(&MakeController{})
}

func (m *MakeController) generateFile(_make, file, desc string) {
	templateFile := m.GetTemplate(_make)

	// 注册lower函数
	//funcMap := template.FuncMap{
	//	"lower": strings.ToLower,
	//	"upper": strings.ToUpper,
	//}

	// 使用Funcs注册函数
	tmpl, err := template.New(filepath.Base(templateFile)). /*Funcs(funcMap).*/ ParseFiles(templateFile)
	if err != nil {
		flag.Errorf("Error parsing template: %s", err.Error())
		os.Exit(1)
	}

	filePath := filepath.ToSlash(file)
	// 去掉/
	filePath = strings.TrimPrefix(filePath, "/")
	// 去掉.go
	filePath = strings.TrimSuffix(filePath, ".go")
	// 文件名
	baseName := path.Base(filePath)
	// 包名
	dir := path.Dir(filePath)

	packageName := path.Base(dir)
	if packageName == "." {
		packageName = ""
	}

	// 路由路径
	routePath := filePath
	// 去掉控制器目录前缀
	routePath = strings.TrimPrefix(routePath, "app/controller/")
	routePath = strings.TrimSuffix(routePath, baseName)
	routePath = routePath + lo.KebabCase(baseName)

	data := struct {
		Package     string
		Name        string
		Description string
		RoutePath   string
	}{
		Package:     packageName,
		Name:        lo.PascalCase(baseName),
		Description: desc,
		RoutePath:   "/api/" + routePath,
	}

	f := m.CheckDirAndFile(file)
	if f == nil {
		return
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		flag.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("控制器文件: %s 生成成功!", file)
}
