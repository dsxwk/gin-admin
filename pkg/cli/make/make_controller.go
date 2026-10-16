package make

import (
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type MakeController struct {
	base.Command
}

// NewMakeController 创建命令
func NewMakeController() *MakeController {
	return &MakeController{
		Command: base.NewCommand(base.Config{
			Name:        "make:controller",
			Description: "控制器创建",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "f", Long: "file"},
					Desc:     "文件路径, 如: v1/user",
					Required: true,
				},
				{
					Flag:     base.Flag{Short: "d", Long: "desc", Default: "用户"},
					Desc:     "描述, 如: 用户",
					Required: false,
				},
			},
		}),
	}
}

func (m *MakeController) Handle(values map[string]string) error {
	_make := strings.TrimPrefix(m.Name(), "make:")
	f := m.GetMakeFile(values["file"], _make)
	m.generateFile(_make, f, values["desc"])
	return nil
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
		terminal.Errorf("Error parsing template: %s", err.Error())
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
		terminal.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	terminal.Successf("控制器文件: %s 生成成功!", file)
}
