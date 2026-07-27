package config

type OperatorRecord struct {
	Enable bool `mapstructure:"enable" yaml:"enable"` // 是否启用操作记录
}
