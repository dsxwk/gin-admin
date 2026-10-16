package console

import (
	"gin/pkg/cli/base"
)

// All 获取全部控制台命令
func All() []base.Handler {
	return []base.Handler{
		NewDemoCommand(),
	}
}
