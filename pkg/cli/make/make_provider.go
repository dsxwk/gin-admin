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

type MakeProvider struct {
	base.BaseCommand
}

func (m *MakeProvider) Name() string {
	return "make:provider"
}

func (m *MakeProvider) Description() string {
	return "创建服务提供者"
}

func (m *MakeProvider) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{
				Short: "f",
				Long:  "file",
			},
			"文件路径, 如: rate_limit",
			true,
		},
		{
			base.Flag{
				Short:   "d",
				Long:    "desc",
				Default: "provider",
			},
			"服务提供者描述, 如: 限流",
			false,
		},
		{
			base.Flag{
				Short:   "D",
				Long:    "deps",
				Default: "config,log",
			},
			"依赖的服务, 多个用逗号分隔, 如: config,log,cache",
			false,
		},
		{
			base.Flag{
				Short: "r",
				Long:  "runner",
			},
			"是否包含后台任务, true/false",
			false,
		},
	}
}

func (m *MakeProvider) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())

	_make := strings.TrimPrefix(m.Name(), "make:")
	fileName := values["file"]

	file := m.GetMakeFile(fileName, _make)
	m.generateProvider(_make, file, values["desc"], values["deps"], values["runner"])
}

func init() {
	cli.Register(&MakeProvider{})
}

func (m *MakeProvider) generateProvider(_make, file, providerDesc, deps, hasRunner string) {
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

	// 从文件名提取提供者名称
	baseName := strings.TrimSuffix(filepath.Base(file), ".go")

	// 转换为大驼峰作为提供者名称
	providerName := pkg.ToUpperCamel(strings.ReplaceAll(baseName, "_", " "))
	providerName = strings.ReplaceAll(providerName, " ", "")

	// 生成小驼峰变量名
	providerVar := pkg.SnakeToCamel(providerName)

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
		ProviderVar:  pkg.CamelToSnake(providerVar),
		Desc:         providerDesc,
		Deps:         dependencyList,
		HasRunner:    withRunner,
	}

	err = tmpl.Execute(f, data)
	if err != nil {
		flag.Errorf("Error executing template: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("服务提供者文件: " + file + " 生成成功!")
}
