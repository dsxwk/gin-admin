package orm

import (
	"gin/pkg"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const startTimeKey = "gorm_start_time"

const (
	defaultMysqlTimeout      = 5 * time.Second
	defaultMysqlReadTimeout  = 5 * time.Second
	defaultMysqlWriteTimeout = 5 * time.Second
)

func openMysql() (db *gorm.DB, err error) {
	db, err = gorm.Open(mysql.Open(getMysqlDsn()), &gorm.Config{
		NamingStrategy: configNaming(),
		Logger:         gormLogger(),
	})

	return db, err
}

// getMysqlDsn 获取数据库dns
func getMysqlDsn() string {
	timeout := conf.Databases.Mysql.Timeout
	if timeout <= 0 {
		timeout = defaultMysqlTimeout
	}

	readTimeout := conf.Databases.Mysql.ReadTimeout
	if readTimeout <= 0 {
		readTimeout = defaultMysqlReadTimeout
	}

	writeTimeout := conf.Databases.Mysql.WriteTimeout
	if writeTimeout <= 0 {
		writeTimeout = defaultMysqlWriteTimeout
	}

	options := pkg.Sprintf(
		"charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai&timeout=%s&readTimeout=%s&writeTimeout=%s&multiStatements=true",
		timeout, readTimeout, writeTimeout,
	)

	return pkg.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		conf.Databases.Mysql.Username, conf.Databases.Mysql.Password, conf.Databases.Mysql.Host, conf.Databases.Mysql.Port, conf.Databases.Mysql.Database, options,
	)
}
