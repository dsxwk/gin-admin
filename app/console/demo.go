package console

import (
	"gin/pkg/cli/base"
	"gin/pkg/terminal"
)

// DemoCommand 演示命令
type DemoCommand struct {
	base.Command
}

// NewDemoCommand 创建演示命令
func NewDemoCommand() *DemoCommand {
	return &DemoCommand{
		Command: base.NewCommand(base.Config{
			Name:        "demo:test",
			Description: "演示控制台命令",
			Options: []base.CommandOption{
				{
					Flag:     base.Flag{Short: "a", Long: "args"},
					Desc:     "演示参数, 如: arg1",
					Required: true,
				},
			},
		}),
	}
}

// Handle 执行命令
func (m *DemoCommand) Handle(values map[string]string) error {
	terminal.Successf("演示参数:%s", values["args"])
	return nil
}
