package config

// App 应用
type App struct {
	Name     string `mapstructure:"name" yaml:"name"`
	Mode     string `mapstructure:"mode" yaml:"mode"`
	Port     int64  `mapstructure:"port" yaml:"port"`
	Timezone string `mapstructure:"timezone" yaml:"timezone"`
	Proxies  string `mapstructure:"proxies" yaml:"proxies"`
	Env      string `mapstructure:"env" yaml:"env"`
}
