package facade

import (
	"gin/config"
)

// Config 配置门面-配置访问统一入口
var Config = &configFacade{}

type configFacade struct{}

// Get 获取配置实例
func (c *configFacade) Get() *config.Config {
	if cfg := Get("config"); cfg != nil {
		return cfg.(*config.Config)
	}
	return config.NewConfig()
}

// GetString 获取字符串配置
func (c *configFacade) GetString(key string) string {
	return config.GetString(key)
}

// GetInt 获取整数配置
func (c *configFacade) GetInt(key string) int {
	return config.GetInt(key)
}

// GetBool 获取布尔配置
func (c *configFacade) GetBool(key string) bool {
	return config.GetBool(key)
}
