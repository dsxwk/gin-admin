package config

// Cors 跨域
type Cors struct {
	Enabled   bool         `mapstructure:"enabled" yaml:"enabled"`
	Mode      string       `mapstructure:"mode" yaml:"mode"`
	AllowAll  *CorsConfig  `mapstructure:"allow-all" yaml:"allow-all"`
	WhiteList []CorsConfig `mapstructure:"whitelist" yaml:"whitelist"`
}

// CorsConfig 跨域配置
type CorsConfig struct {
	AllowOrigin      string `mapstructure:"allow-origin" yaml:"allow-origin"`
	AllowHeaders     string `mapstructure:"allow-headers" yaml:"allow-headers"`
	ExposeHeaders    string `mapstructure:"expose-headers" yaml:"expose-headers"`
	AllowMethods     string `mapstructure:"allow-methods" yaml:"allow-methods"`
	AllowCredentials string `mapstructure:"allow-credentials" yaml:"allow-credentials"`
}

// GetConfig 根据模式获取当前生效的跨域配置
func (c *Cors) GetConfig(origin string) *CorsConfig {
	if !c.Enabled {
		return nil
	}

	switch c.Mode {
	case "allow-all":
		return c.AllowAll
	case "whitelist":
		for _, item := range c.WhiteList {
			if item.AllowOrigin == origin || item.AllowOrigin == "*" {
				return &item
			}
		}
		return nil
	default:
		return nil
	}
}
