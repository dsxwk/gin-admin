package config

type OperatorRecord struct {
	Enable bool `mapstructure:"enabled" yaml:"enabled"` // 是否启用操作记录
}
