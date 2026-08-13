package config

// Agent AI Agent配置
type Agent struct {
	Enabled     bool                     `mapstructure:"enabled" yaml:"enabled"`         // 是否启用Agent
	Default     string                   `mapstructure:"default" yaml:"default"`         // 默认提供商
	MaxTokens   int                      `mapstructure:"max-tokens" yaml:"max-tokens"`   // 单次最大token
	Temperature float64                  `mapstructure:"temperature" yaml:"temperature"` // 温度参数
	Providers   map[string]AgentProvider `mapstructure:"providers" yaml:"providers"`     // 模型提供商列表
}

// AgentProvider 模型提供商配置
type AgentProvider struct {
	ApiKey  string `mapstructure:"api-key" yaml:"api-key"`   // APIKey
	BaseUrl string `mapstructure:"base-url" yaml:"base-url"` // 接口地址
	Model   string `mapstructure:"model" yaml:"model"`       // 模型名称
}
