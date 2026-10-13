package facade

import (
	"gin/pkg/serviceprovider/logger"
)

// Log 日志门面方法
func Log() *logger.Logger {
	return App().Log()
}
