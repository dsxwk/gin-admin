package make

import (
	"gin/pkg/cli"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type MakeRouter struct {
	base.Command
}

// NewMakeRouter 创建命令
func NewMakeRouter() *MakeRouter {
	return &MakeRouter{
		Command: base.NewCommand(base.Config{
			Name:        "make:router",
			Description: "路由创建",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "f", Long: "file"},
					Desc:     "文件路径, 如: user",
					Required: true,
				},
				{
					Flag:     base.Flag{Short: "d", Long: "desc", Default: ""},
					Desc:     "路由描述, 如: 用户路由",
					Required: false,
				},
			},
		}),
	}
}

func (m *MakeRouter) Handle(values map[string]string) error {
	_make := strings.TrimPrefix(m.Name(), "make:")
	f := m.GetMakeFile(values["file"], _make)
	m.generateFile(_make, f, values["desc"])
	return nil
}

func (m *MakeRouter) generateFile(_make, file, desc string) {
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

	name := strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))

	if desc == "" {
		desc = lo.PascalCase(name)
	}

	data := struct {
		Package       string // 提取的包名
		Name          string // 模块名称(首字母大写)
		NameLower     string // 名称(首字母小写)
		NameKebabCase string // 名称(中划线)
		Description   string // 如果为空,使用默认值
	}{
		Package:       packageName,
		Name:          lo.PascalCase(name),
		NameLower:     lo.CamelCase(name),
		NameKebabCase: lo.KebabCase(name),
		Description:   desc,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		terminal.Successf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	registryFile := filepath.Join("router", "registry.go")
	qualifier := ""
	routerDir := filepath.ToSlash(filepath.Clean(filepath.Dir(file)))
	if routerDir != "router" {
		alias, importErr := cli.AddImport(registryFile, cli.ImportPath(file))
		if importErr != nil {
			terminal.Errorf("自动添加路由导入失败: %s", importErr.Error())
			os.Exit(1)
		}
		qualifier = alias + "."
	}

	item := "&" + qualifier + data.Name + "Router{}"
	if err = cli.AddItem(registryFile, "return []route.Router{", item); err != nil {
		terminal.Errorf("自动注册路由失败: %s", err.Error())
		os.Exit(1)
	}

	terminal.Successf("路由文件: " + file + " 生成成功!")
}
