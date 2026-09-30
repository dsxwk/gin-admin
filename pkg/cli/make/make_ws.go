package make

import (
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/samber/lo"
)

// MakeWs WebSocket消息处理器创建命令
type MakeWs struct {
	base.BaseCommand
}

// Name 命令名称
func (m *MakeWs) Name() string {
	return "make:ws"
}

// Description 命令描述
func (m *MakeWs) Description() string {
	return "创建ws消息处理器"
}

// Help 命令选项
func (m *MakeWs) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			Flag:     base.Flag{Short: "f", Long: "file"},
			Desc:     "文件路径, 如: notice",
			Required: true,
		},
		{
			Flag:     base.Flag{Short: "t", Long: "type"},
			Desc:     "消息类型, 默认取文件名的蛇形命名",
			Required: false,
		},
		{
			Flag:     base.Flag{Short: "d", Long: "desc", Default: "消息处理"},
			Desc:     "处理器描述, 如: 通知消息",
			Required: false,
		},
	}
}

// Execute 执行命令
func (m *MakeWs) Execute(values map[string]string) {
	_make := strings.TrimPrefix(m.Name(), "make:")
	file := m.GetMakeFile(values["file"], _make)
	m.generateFile(file, values["file"], values["type"], values["desc"])
}

func init() {
	cli.Register(&MakeWs{})
}

// generateFile 生成WebSocket消息处理器
func (m *MakeWs) generateFile(file, rawName, messageType, description string) {
	templateFile := m.GetTemplate("ws")
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		flag.Errorf("Error parsing ws template: %s", err.Error())
		os.Exit(1)
	}

	baseName := strings.TrimSuffix(filepath.Base(filepath.ToSlash(rawName)), ".go")
	name := lo.PascalCase(baseName)
	if messageType == "" {
		messageType = lo.SnakeCase(baseName)
	}

	data := struct {
		Package     string
		Name        string
		MessageType string
		Description string
	}{
		Package:     filepath.Base(filepath.Dir(file)),
		Name:        name,
		MessageType: messageType,
		Description: description,
	}

	output := m.CheckDirAndFile(file)
	if output == nil {
		return
	}
	defer func() { _ = output.Close() }()

	if err = tmpl.Execute(output, data); err != nil {
		flag.Errorf("Error executing ws template: %s", err.Error())
		os.Exit(1)
	}

	registryFile := filepath.Join("websocket", "registry.go")
	qualifier := ""
	wsDir := filepath.ToSlash(filepath.Clean(filepath.Dir(file)))
	if wsDir != "websocket" {
		alias, importErr := addRegistryImport(registryFile, registryImportPath(file))
		if importErr != nil {
			flag.Errorf("自动添加ws处理器导入失败: %s", importErr.Error())
			return
		}
		qualifier = alias + "."
	}

	item := "&" + qualifier + name + "Handler{}"
	if err = addRegistryItem(registryFile, "return []ws.Handler{", item); err != nil {
		flag.Errorf("自动注册ws处理器失败: %s", err.Error())
		return
	}

	flag.Successf("ws处理器文件: %s 生成成功!", file)
}
