package make

import (
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
)

type MakeConsole struct {
	base.BaseCommand
}

func (m *MakeConsole) Name() string {
	return "make:console"
}

func (m *MakeConsole) Description() string {
	return "控制台命令创建"
}

func (m *MakeConsole) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "n",
				Long:  "name",
			},
			"命令名称, 如: demo:test",
			true,
		},
		{
			base.Flag{
				Short: "p",
				Long:  "path",
			},
			"文件目录, 如: app/console/cronjob, 默认 app/console",
			false,
		},
		{
			base.Flag{
				Short:   "d",
				Long:    "desc",
				Default: "console-desc",
			},
			"描述, 如: console-desc",
			false,
		},
	}
}

func (m *MakeConsole) Execute(values map[string]string) {
	file := m.resolveFile(values)
	m.generateFile(file, values["name"], values["desc"])
}

func (m *MakeConsole) generateFile(file, name, desc string) {
	templateFile := m.GetTemplate("console")
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
	defer func() { _ = f.Close() }()

	baseName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(filepath.Base(file)))
	name = strings.TrimSpace(name)

	if strings.TrimSpace(desc) == "" {
		desc = "console-desc"
	}

	data := struct {
		Package     string // 提取的包名
		Name        string // 模块名称(首字母大写)
		Command     string // 如果为空,使用默认值
		Description string // 如果为空,使用默认值
	}{
		Package:     packageName,
		Name:        lo.PascalCase(baseName),
		Command:     name,
		Description: desc,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		flag.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	commandDir := filepath.ToSlash(filepath.Clean(filepath.Dir(file)))
	registryFile := filepath.Join("app", "console", "registry.go")
	qualifier := ""
	if commandDir != "app/console" {
		alias, importErr := cli.AddImport(registryFile, cli.ImportPath(file))
		if importErr != nil {
			flag.Errorf("自动添加控制台命令导入失败: %s", importErr.Error())
			os.Exit(1)
		}
		qualifier = alias + "."
	}

	item := "&" + qualifier + lo.PascalCase(baseName) + "Command{}"
	if err = cli.AddItem(registryFile, "return []base.Command{", item); err != nil {
		flag.Errorf("自动注册控制台命令失败: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("控制台命令文件: " + file + " 生成成功!")
}

// resolveFile 解析控制台命令文件路径
func (m *MakeConsole) resolveFile(values map[string]string) string {
	baseName := lo.SnakeCase(strings.TrimSpace(values["name"]))

	dir := strings.TrimSpace(values["path"])
	if dir == "" {
		return filepath.Join("app", "console", baseName+".go")
	}

	dir = filepath.FromSlash(strings.TrimPrefix(filepath.ToSlash(dir), "/"))
	if strings.HasSuffix(strings.ToLower(dir), ".go") {
		return dir
	}
	return filepath.Join(dir, baseName+".go")
}
