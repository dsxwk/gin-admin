package command

import (
	"gin/common/base"
	"gin/pkg/cli"
)

type DemoCommand struct {
	base.BaseCommand
}

func (m *DemoCommand) Name() string {
	return "demo:command"
}

func (m *DemoCommand) Description() string {
	return "test-demo"
}

func (m *DemoCommand) Help() []base.CommandOption {
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

func (m *DemoCommand) Execute(args []string) {
	_ = m.ParseFlags(m.Name(), args, m.Help())
}

func init() {
	cli.Register(&DemoCommand{})
}
