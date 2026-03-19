package config

// Cors 跨域
type Cors struct {
	Enabled          bool   `mapstructure:"enabled" yaml:"enabled"`
	AllowOrigin      string `mapstructure:"allow-origin" yaml:"allow-origin"`
	AllowHeaders     string `mapstructure:"allow-headers" yaml:"allow-headers"`
	ExposeHeaders    string `mapstructure:"expose-headers" yaml:"expose-headers"`
	AllowMethods     string `mapstructure:"allow-methods" yaml:"allow-methods"`
	AllowCredentials string `mapstructure:"allow-credentials" yaml:"allow-credentials"`
}
