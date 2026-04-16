package orm

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
		conf.Databases.Mysql.Username, conf.Databases.Mysql.Password, conf.Databases.Mysql.Host, conf.Databases.Mysql.Port, conf.Databases.Mysql.Database,
	)
}
