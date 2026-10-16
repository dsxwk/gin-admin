package {{.Package}}

import (
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
)

// {{.Name}}Command {{.Description}}
type {{.Name}}Command struct {
	base.Command
}

// New{{.Name}}Command 创建{{.Description}}命令
func New{{.Name}}Command() *{{.Name}}Command {
	return &{{.Name}}Command{
		Command: base.NewCommand(base.Config{
			Name: {{if eq .Command ""}}"{{.Name}}-command"{{else}}"{{.Command}}"{{end}},
			Description: "{{.Description}}",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "a", Long: "args"},
					Desc:     "示例参数, 如: arg1",
					Required: true,
				},
			},
		}),
	}
}

// Handle 执行命令
func (m *{{.Name}}Command) Handle(values map[string]string) error {
	terminal.Successf("执行命令: %s %s", m.Name(), m.FormatArgs(values))
	return nil
}
