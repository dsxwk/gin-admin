package facade

import (
	"gin/pkg/container"

	"gorm.io/gorm"
)

// DB 数据库门面-数据库访问统一入口
func DB(conn ...string) *gorm.DB {
	manager := container.Default().DB()
	if manager == nil {
		return nil
	}

	name := ""
	if len(conn) > 0 && conn[0] != "" {
		name = conn[0]
	}

	db := manager.Connection(name)
	conf := Config()
	if conf == nil {
		return nil
	}
	if conf.Databases.DisableSoftDelete {
		db = db.Unscoped()
	}
	return db
}

// ResetDB 重置数据库连接,关闭旧连接并重建
func ResetDB(db *gorm.DB) *gorm.DB {
	conf := Config()
	manager := container.Default().DB()
	if manager == nil || conf == nil {
		return nil
	}

	newDB := manager.Reset(conf.Databases.Driver)
	if conf.Databases.DisableSoftDelete {
		newDB = newDB.Unscoped()
	}
	return newDB
}
