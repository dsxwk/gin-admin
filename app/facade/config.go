package facade

import (
	"gin/config"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
)

// Config 配置门面方法
func Config() *config.Config {
	return container.Default().Get[*config.Config](serviceprovider.ServiceConfig)
}
