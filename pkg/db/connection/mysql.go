package connection

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
		conf.Mysql.Username, conf.Mysql.Password, conf.Mysql.Host, conf.Mysql.Port, conf.Mysql.Database,
	)
}
