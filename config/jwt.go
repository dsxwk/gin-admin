package config

// Jwt token
type Jwt struct {
	Key        string `mapstructure:"key" yaml:"key"`
	Exp        int64  `mapstructure:"exp" yaml:"exp"`
	RefreshExp int64  `mapstructure:"refresh-exp" yaml:"refresh-exp"`
}
