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

// MakeWs WebSocket消息处理器创建命令
type MakeWs struct {
	base.Command
}

// NewMakeWs 创建命令
func NewMakeWs() *MakeWs {
	return &MakeWs{
		Command: base.NewCommand(base.Config{
			Name:        "make:ws",
			Description: "创建ws消息处理器",
			Options: []base.CommandOption{
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
			},
		}),
	}
}

// Handle 执行命令
func (m *MakeWs) Handle(values map[string]string) error {
	_make := strings.TrimPrefix(m.Name(), "make:")
	file := m.GetMakeFile(values["file"], _make)
	m.generateFile(file, values["file"], values["type"], values["desc"])
	return nil
}

// generateFile 生成WebSocket消息处理器
func (m *MakeWs) generateFile(file, rawName, messageType, description string) {
	templateFile := m.GetTemplate("ws")
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		terminal.Errorf("Error parsing ws template: %s", err.Error())
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
		terminal.Errorf("Error executing ws template: %s", err.Error())
		os.Exit(1)
	}

	registryFile := filepath.Join("websocket", "registry.go")
	qualifier := ""
	wsDir := filepath.ToSlash(filepath.Clean(filepath.Dir(file)))
	if wsDir != "websocket" {
		alias, importErr := cli.AddImport(registryFile, cli.ImportPath(file))
		if importErr != nil {
			terminal.Errorf("自动添加ws处理器导入失败: %s", importErr.Error())
			return
		}
		qualifier = alias + "."
	}

	item := "&" + qualifier + name + "Handler{}"
	if err = cli.AddItem(registryFile, "return []ws.Handler{", item); err != nil {
		terminal.Errorf("自动注册ws处理器失败: %s", err.Error())
		return
	}

	terminal.Successf("ws处理器文件: %s 生成成功!", file)
}
