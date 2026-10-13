package facade

import (
	"gin/pkg/serviceprovider/debugger"
)

// Debugger 获取调试器
func Debugger() *debugger.Debugger {
	return App().Debugger()
}
