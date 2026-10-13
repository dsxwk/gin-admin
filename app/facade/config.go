package facade

import (
	"gin/config"
)

// Config 配置门面方法
func Config() *config.Config {
	return App().Config()
}
