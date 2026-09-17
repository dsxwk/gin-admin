package facade

import (
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/logger"
)

// Log 日志门面方法
func Log() *logger.Logger {
	return container.Default().Get[*logger.Logger](serviceprovider.ServiceLog)
}
