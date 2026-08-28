package config

// Grpc grpc服务配置
type Grpc struct {
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"` // 是否启用grpc服务
	Host    string `mapstructure:"host" yaml:"host"`       // 监听地址
	Port    int    `mapstructure:"port" yaml:"port"`       // 监听端口
}
