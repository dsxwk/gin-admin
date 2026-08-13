package config

// Mcp MCP服务配置
type Mcp struct {
	Enabled bool    `mapstructure:"enabled" yaml:"enabled"` // 是否启用MCP服务
	Path    string  `mapstructure:"path" yaml:"path"`       // MCP路由前缀
	Auth    McpAuth `mapstructure:"auth" yaml:"auth"`       // 认证配置
}

// McpAuth MCP认证配置
type McpAuth struct {
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"` // 是否启用认证
	Token   string `mapstructure:"token" yaml:"token"`     // BearerToken
}
