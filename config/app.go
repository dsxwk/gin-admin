package config

import "time"

// App 应用
type App struct {
	CliVersion string        `mapstructure:"cli-version" yaml:"cli-version"`
	Env        string        `mapstructure:"env" yaml:"env"`
	Name       string        `mapstructure:"name" yaml:"name"`
	Mode       string        `mapstructure:"mode" yaml:"mode"`
	Port       int64         `mapstructure:"port" yaml:"port"`
	Timezone   string        `mapstructure:"timezone" yaml:"timezone"`
	Proxies    string        `mapstructure:"proxies" yaml:"proxies"`
	Timeout    time.Duration `mapstructure:"timeout" yaml:"timeout"`
}
