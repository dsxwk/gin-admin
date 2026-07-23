package config

type Queue struct {
	Connection string   `mapstructure:"connection" yaml:"connection"`
	Redis      Redis    `mapstructure:"redis" yaml:"redis"`
	Kafka      Kafka    `mapstructure:"kafka" yaml:"kafka"`
	Rabbitmq   Rabbitmq `mapstructure:"rabbitmq" yaml:"rabbitmq"`
}

type Kafka struct {
	Enabled bool     `mapstructure:"enabled" yaml:"enabled"` // 鏄惁鍚敤
	Brokers []string `mapstructure:"brokers" yaml:"brokers"`
}

type Rabbitmq struct {
	Enabled bool   `mapstructure:"enabled" yaml:"enabled"` // 鏄惁鍚敤
	Url     string `mapstructure:"url" yaml:"url"`
}
