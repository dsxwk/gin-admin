package config

type Kafka struct {
	Enabled bool     `mapstructure:"enabled" yaml:"enabled"` // 是否启用
	Brokers []string `mapstructure:"brokers" yaml:"brokers"`
}

type Rabbitmq struct {
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"` // 是否启用
	Url     string `mapstructure:"url" yaml:"url"`
}
