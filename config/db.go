package config

import "time"

// Databases 数据库
type Databases struct {
	DbConnection      string        `mapstructure:"db-connection" yaml:"db-connection"`             // 默认数据库
	SlowQueryDuration time.Duration `mapstructure:"slow-query-duration" yaml:"slow-query-duration"` // 慢查询的时间(ms) 超过这个时间会记录到日志中
}

// Mysql 数据库
type Mysql struct {
	Driver            string        `mapstructure:"driver" yaml:"driver"`
	Host              string        `mapstructure:"host" yaml:"host"`
	Port              string        `mapstructure:"port" yaml:"port"`
	Database          string        `mapstructure:"database" yaml:"database"`
	Username          string        `mapstructure:"username" yaml:"username"`
	Password          string        `mapstructure:"password" yaml:"password"`
	SlowQueryDuration time.Duration `mapstructure:"slow-query-duration" yaml:"slow-query-duration"` // 慢查询的时间(ms) 超过这个时间会记录到日志中
}

// Sqlite 数据库
type Sqlite struct {
	Driver            string        `mapstructure:"driver" yaml:"driver"`
	Path              string        `mapstructure:"path" yaml:"path"`
	SlowQueryDuration time.Duration `mapstructure:"slow-query-duration" yaml:"slow-query-duration"` // 慢查询的时间(ms) 超过这个时间会记录到日志中
}

// Pgsql 数据库
type Pgsql struct {
	Driver            string        `mapstructure:"driver" yaml:"driver"`
	Host              string        `mapstructure:"host" yaml:"host"`
	Port              string        `mapstructure:"port" yaml:"port"`
	Database          string        `mapstructure:"database" yaml:"database"`
	Username          string        `mapstructure:"username" yaml:"username"`
	Password          string        `mapstructure:"password" yaml:"password"`
	SlowQueryDuration time.Duration `mapstructure:"slow-query-duration" yaml:"slow-query-duration"` // 慢查询的时间(ms) 超过这个时间会记录到日志中
}

// Sqlsrv 数据库
type Sqlsrv struct {
	Driver            string        `mapstructure:"driver" yaml:"driver"`
	Host              string        `mapstructure:"host" yaml:"host"`
	Port              string        `mapstructure:"port" yaml:"port"`
	Database          string        `mapstructure:"database" yaml:"database"`
	Username          string        `mapstructure:"username" yaml:"username"`
	Password          string        `mapstructure:"password" yaml:"password"`
	SlowQueryDuration time.Duration `mapstructure:"slow-query-duration" yaml:"slow-query-duration"` // 慢查询的时间(ms) 超过这个时间会记录到日志中
}
