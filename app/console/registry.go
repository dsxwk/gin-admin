package console

import (
	"gin/common/base"
)

// All 获取全部控制台命令
func All() []base.Command {
	return []base.Command{
		&DemoCommand{},
	}
}
