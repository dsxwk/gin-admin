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

// ResetDB 重置数据库连接,关闭旧连接并重建
func ResetDB(db *gorm.DB) *gorm.DB {
	conf := Config()
	name := conf.Databases.Driver

	mgr := GetManager()
	mgr.mu.Lock()
	delete(mgr.instances, name)
	mgr.mu.Unlock()

	newDB := orm.ResetConnection(name)

	Register[*gorm.DB](name, newDB)

	if conf.Databases.DisableSoftDelete {
		newDB = newDB.Unscoped()
	}
	return newDB
}
