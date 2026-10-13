package console

import (
	"gin/common/base"
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

func (m *DemoCommand) Execute(_ map[string]string) {
	// todo
}
