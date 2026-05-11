package facade

import (
	"gin/config"
)

// Config 配置门面方法
func Config() *config.Config {
	cfg := Get[*config.Config]("config")
	if cfg != nil {
		return cfg
	}
	return config.NewConfig()
}
