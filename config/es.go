package config

import "time"

// Es Elasticsearch配置
type Es struct {
	Enabled   bool          `mapstructure:"enabled" yaml:"enabled"`     // 是否启用ES
	Addresses []string      `mapstructure:"addresses" yaml:"addresses"` // ES地址列表
	Username  string        `mapstructure:"username" yaml:"username"`   // 用户名
	Password  string        `mapstructure:"password" yaml:"password"`   // 密码
	ApiKey    string        `mapstructure:"api-key" yaml:"api-key"`     // APIKey
	Timeout   time.Duration `mapstructure:"timeout" yaml:"timeout"`     // 请求超时
}
