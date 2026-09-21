package facade

import (
	"gin/config"
	"gin/pkg/container"
)

// Config 配置门面方法
func Config() *config.Config {
	return container.Default().Config()
}
