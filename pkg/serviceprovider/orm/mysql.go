package orm

import (
	"fmt"
	"gin/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const (
	defaultMysqlTimeout      = 5 * time.Second
	defaultMysqlReadTimeout  = 5 * time.Second
	defaultMysqlWriteTimeout = 5 * time.Second
)

func openMysql(cfg *config.Config, gormLog gormlogger.Interface) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(getMysqlDsn(cfg)), &gorm.Config{
		NamingStrategy: configNaming(),
		Logger:         gormLog,
	})
}

// getMysqlDsn 获取数据库dns
func getMysqlDsn(cfg *config.Config) string {
	timeout := cfg.Databases.Mysql.Timeout
	if timeout <= 0 {
		timeout = defaultMysqlTimeout
	}

	readTimeout := cfg.Databases.Mysql.ReadTimeout
	if readTimeout <= 0 {
		readTimeout = defaultMysqlReadTimeout
	}

	writeTimeout := cfg.Databases.Mysql.WriteTimeout
	if writeTimeout <= 0 {
		writeTimeout = defaultMysqlWriteTimeout
	}

	options := fmt.Sprintf(
		"charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai&timeout=%s&readTimeout=%s&writeTimeout=%s&multiStatements=true",
		timeout, readTimeout, writeTimeout,
	)

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		cfg.Databases.Mysql.Username, cfg.Databases.Mysql.Password, cfg.Databases.Mysql.Host, cfg.Databases.Mysql.Port, cfg.Databases.Mysql.Database, options,
	)
}
