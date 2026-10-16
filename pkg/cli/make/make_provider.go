package make

import (
	"gin/pkg/cli"
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/samber/lo"
)

type MakeProvider struct {
	base.Command
}

// NewMakeProvider 创建命令
func NewMakeProvider() *MakeProvider {
	return &MakeProvider{
		Command: base.NewCommand(base.Config{
			Name:        "make:provider",
			Description: "创建服务提供者",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "f", Long: "file"},
					Desc:     "文件路径, 如: rate_limit",
					Required: true,
				},
				{
					Flag:     base.Flag{Short: "d", Long: "desc", Default: "provider"},
					Desc:     "服务提供者描述, 如: 限流",
					Required: false,
				},
				{
					Flag:     base.Flag{Short: "D", Long: "deps", Default: "config,log"},
					Desc:     "依赖的服务, 多个用逗号分隔, 如: config,log,cache",
					Required: false,
				},
				{
					Flag:     base.Flag{Short: "r", Long: "runner"},
					Desc:     "是否包含后台任务, true/false",
					Required: false,
				},
			},
		}),
	}
}

func (m *MakeProvider) Handle(values map[string]string) error {

	_make := strings.TrimPrefix(m.Name(), "make:")
	fileName := values["file"]

	file := m.GetMakeFile(fileName, _make)
	m.generateProvider(_make, file, values["desc"], values["deps"], values["runner"])
	return nil
}

func (m *MakeProvider) generateProvider(_make, file, providerDesc, deps, hasRunner string) {
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

	// 从文件名提取提供者名称
	baseName := strings.TrimSuffix(filepath.Base(file), ".go")

	// 转换为大驼峰作为提供者名称
	providerName := lo.PascalCase(strings.ReplaceAll(baseName, "_", " "))
	providerName = strings.ReplaceAll(providerName, " ", "")

	// 生成小驼峰变量名
	providerVar := lo.CamelCase(providerName)

	// 处理依赖
	dependencyList := strings.Split(deps, ",")
	for i := range dependencyList {
		dependencyList[i] = strings.TrimSpace(dependencyList[i])
	}

	// 是否有Runner
	withRunner := hasRunner == "true" || hasRunner == "1" || hasRunner == "yes"

	data := struct {
		Package      string   // 包名(provider)
		ProviderName string   // 提供者名称(大驼峰)
		ProviderVar  string   // 提供者变量名(小写驼峰)
		Desc         string   // 提供者描述
		Deps         []string // 依赖列表
		HasRunner    bool     // 是否有后台任务
	}{
		Package:      "provider",
		ProviderName: providerName,
		ProviderVar:  lo.SnakeCase(providerVar),
		Desc:         providerDesc,
		Deps:         dependencyList,
		HasRunner:    withRunner,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		terminal.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	providersFile := filepath.Join("app", "provider", "registry.go")
	qualifier := ""
	providerDir := filepath.ToSlash(filepath.Clean(filepath.Dir(file)))
	if providerDir != "app/provider" {
		alias, importErr := cli.AddImport(providersFile, cli.ImportPath(file))
		if importErr != nil {
			terminal.Errorf("自动添加服务提供者导入失败: %s", importErr.Error())
			return
		}
		qualifier = alias + "."
	}
	providerItem := "&" + qualifier + providerName + "Provider{}"
	if err = cli.AddItem(providersFile, "return []serviceprovider.ServiceProvider{", providerItem); err != nil {
		terminal.Errorf("自动注册服务提供者失败: %s", err.Error())
		return
	}

	terminal.Successf("服务提供者文件: " + file + " 生成成功!")
}
