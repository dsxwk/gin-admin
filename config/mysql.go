package config

import (
	"gin/pkg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const startTimeKey = "gorm_start_time"

func openMysql() (db *gorm.DB, err error) {
	db, err = gorm.Open(mysql.Open(getMysqlDsn()), &gorm.Config{
		NamingStrategy: configNaming(),
		Logger:         gormLogger(),
	})

	return db, err
}

// getMysqlDsn 获取数据库dns
func getMysqlDsn() string {
	return pkg.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		Conf.Mysql.Username, Conf.Mysql.Password, Conf.Mysql.Host, Conf.Mysql.Port, Conf.Mysql.Database,
	)
}
