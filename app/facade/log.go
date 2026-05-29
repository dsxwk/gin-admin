package facade

import (
	"gin/pkg/serviceprovider/logger"
)

// Log 日志门面方法
func Log() *logger.Logger {
	log := Get[*logger.Logger]("log")
	if log != nil {
		return log
	}
	return logger.NewLogger(Config())
}
