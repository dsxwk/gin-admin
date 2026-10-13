package {{.Package}}

import (
	"gin/common/base"
	"github.com/fatih/color"
)

type {{.Name}}Command struct {
	base.BaseCommand
}

func (m *{{.Name}}Command) Name() string {
    {{- if eq .Command "" }}
    return "{{.Name}}-command"
    {{- else }}
    return "{{.Command}}"
    {{- end }}
}

func (m *{{.Name}}Command) Description() string {
	return "{{.Description}}"
}

func (m *{{.Name}}Command) Help() []base.CommandOption {
	return []base.CommandOption{
        {
            base.Flag{
                Short: "a",
                Long:  "args",
            },
            "示例参数, 如: arg1",
            true,
        },
    }
}

func (m *{{.Name}}Command) Execute(values map[string]string) {
    _ = values
    color.Green("执行命令: %s %s", m.Name(), m.FormatArgs(values))
}
