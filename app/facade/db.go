package facade

import (
	"gin/pkg/serviceprovider/orm"
	"gorm.io/gorm"
)

// DB 数据库门面-数据库访问统一入口
func DB(conn ...string) *gorm.DB {
	conf := Config()
	name := conf.Databases.Driver
	if len(conn) > 0 && conn[0] != "" {
		name = conn[0]
	}

	db := Get[*gorm.DB](name)
	if db != nil {
		if conf.Databases.DisableSoftDelete {
			db = db.Unscoped()
		}
		return db
	}

	db = orm.Connection(name, conf)
	if conf.Databases.DisableSoftDelete {
		db = db.Unscoped()
	}
	return db
}
