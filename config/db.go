package config

import "time"

// Databases 数据库
type Databases struct {
	Default           string        `mapstructure:"default" yaml:"default"`                         // 默认数据库
	SlowQueryDuration time.Duration `mapstructure:"slow-query-duration" yaml:"slow-query-duration"` // 慢查询的时间,超过这个时间会记录到日志中
	Mysql             Mysql         `mapstructure:"mysql" yaml:"mysql"`                             // mysql
	Sqlite            Sqlite        `mapstructure:"sqlite" yaml:"sqlite"`                           // sqlite
	Pgsql             Pgsql         `mapstructure:"pgsql" yaml:"pgsql"`                             // pgsql
	Sqlsrv            Sqlsrv        `mapstructure:"sqlsrv" yaml:"sqlsrv"`                           // sqlsrv
}

// Mysql 数据库
type Mysql struct {
	Driver   string `mapstructure:"driver" yaml:"driver"`
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	Database string `mapstructure:"database" yaml:"database"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
}

// Sqlite 数据库
type Sqlite struct {
	Driver string `mapstructure:"driver" yaml:"driver"`
	Path   string `mapstructure:"path" yaml:"path"`
}

// Pgsql 数据库
type Pgsql struct {
	Driver   string `mapstructure:"driver" yaml:"driver"`
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	Database string `mapstructure:"database" yaml:"database"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
	SSLMode  string `mapstructure:"ssl-mode" yaml:"ssl-mode"`
}

// Sqlsrv 数据库
type Sqlsrv struct {
	Driver   string `mapstructure:"driver" yaml:"driver"`
	Host     string `mapstructure:"host" yaml:"host"`
	Port     string `mapstructure:"port" yaml:"port"`
	Database string `mapstructure:"database" yaml:"database"`
	Username string `mapstructure:"username" yaml:"username"`
	Password string `mapstructure:"password" yaml:"password"`
}
