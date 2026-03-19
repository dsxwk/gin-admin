package config

// Log 日志
type Log struct {
	Access     bool   `mapstructure:"access" yaml:"access"`           // 是否记录访问日志
	MaxSize    int    `mapstructure:"max-size" yaml:"max-size"`       // 单个日志文件大小（MB）
	MaxBackups int    `mapstructure:"max-backups" yaml:"max-backups"` // 最多保留的旧日志文件数
	MaxDay     int    `mapstructure:"max-day" yaml:"max-day"`         // 保留的最大天数
	Level      string `mapstructure:"level" yaml:"level"`             // 日志级别
}
