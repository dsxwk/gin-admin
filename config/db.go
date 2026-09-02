package config

import "time"

// Databases 数据库
type Databases struct {
	Driver            string        `mapstructure:"driver" yaml:"driver"`                           // 默认数据库
	DisableSoftDelete bool          `mapstructure:"disable-soft-delete" yaml:"disable-soft-delete"` // 禁用软删除
	SlowQueryDuration time.Duration `mapstructure:"slow-query-duration" yaml:"slow-query-duration"` // 慢查询的时间,超过这个时间会记录到日志中
	Pool              DatabasePool  `mapstructure:"pool" yaml:"pool"`                               // 连接池
	Mysql             Mysql         `mapstructure:"mysql" yaml:"mysql"`                             // mysql
	Sqlite            Sqlite        `mapstructure:"sqlite" yaml:"sqlite"`                           // sqlite
	Pgsql             Pgsql         `mapstructure:"pgsql" yaml:"pgsql"`                             // pgsql
	Sqlsrv            Sqlsrv        `mapstructure:"sqlsrv" yaml:"sqlsrv"`                           // sqlsrv
	Oracle            Oracle        `mapstructure:"oracle" yaml:"oracle"`                           // oracle
}

// DatabasePool 数据库连接池
type DatabasePool struct {
	MaxIdleConns    int           `mapstructure:"max-idle-conns" yaml:"max-idle-conns"`         // 最大空闲连接数
	MaxOpenConns    int           `mapstructure:"max-open-conns" yaml:"max-open-conns"`         // 最大打开连接数
	ConnMaxLifetime time.Duration `mapstructure:"conn-max-lifetime" yaml:"conn-max-lifetime"`   // 连接最大存活时间
	ConnMaxIdleTime time.Duration `mapstructure:"conn-max-idle-time" yaml:"conn-max-idle-time"` // 连接最大空闲时间
}

// Mysql 数据库
type Mysql struct {
	Driver       string        `mapstructure:"driver" yaml:"driver"`
	Host         string        `mapstructure:"host" yaml:"host"`
	Port         string        `mapstructure:"port" yaml:"port"`
	Database     string        `mapstructure:"database" yaml:"database"`
	Username     string        `mapstructure:"username" yaml:"username"`
	Password     string        `mapstructure:"password" yaml:"password"`
	Timeout      time.Duration `mapstructure:"timeout" yaml:"timeout"`             // 连接超时
	ReadTimeout  time.Duration `mapstructure:"read-timeout" yaml:"read-timeout"`   // 读取超时
	WriteTimeout time.Duration `mapstructure:"write-timeout" yaml:"write-timeout"` // 写入超时
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

// Oracle 数据库
type Oracle struct {
	Host        string `yaml:"host"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	ServiceName string `yaml:"service-name"`
	Sid         string `yaml:"sid"` // 可选,与ServiceName二选一
	Port        string `yaml:"port"`
}
