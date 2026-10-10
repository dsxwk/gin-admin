package config

import "time"

// Agent AI Agent配置
type Agent struct {
	Enabled        bool                     `mapstructure:"enabled" yaml:"enabled"`                 // 是否启用Agent
	Default        string                   `mapstructure:"default" yaml:"default"`                 // 默认提供商
	MaxTokens      int                      `mapstructure:"max-tokens" yaml:"max-tokens"`           // 单次最大token
	Temperature    float64                  `mapstructure:"temperature" yaml:"temperature"`         // 温度参数
	RequestTimeout time.Duration            `mapstructure:"request-timeout" yaml:"request-timeout"` // 模型请求超时
	Sse            AgentSse                 `mapstructure:"sse" yaml:"sse"`                         // SSE配置
	Providers      map[string]AgentProvider `mapstructure:"providers" yaml:"providers"`             // 模型提供商列表
}

// AgentSse SSE配置
type AgentSse struct {
	Heartbeat time.Duration `mapstructure:"heartbeat" yaml:"heartbeat"` // 心跳间隔
	Retry     int           `mapstructure:"retry" yaml:"retry"`         // 重连时间(毫秒)
}

// AgentProvider 模型提供商配置
type AgentProvider struct {
	ApiKey  string `mapstructure:"api-key" yaml:"api-key"`   // APIKey
	BaseUrl string `mapstructure:"base-url" yaml:"base-url"` // 接口地址
	Model   string `mapstructure:"model" yaml:"model"`       // 模型名称
}
