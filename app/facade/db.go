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

	return manager.Connection(name)
}

// ResetDB 重置数据库连接,关闭旧连接并重建
func ResetDB(db *gorm.DB) *gorm.DB {
	manager := container.Default().DB()
	if manager == nil {
		return nil
	}

	return manager.Reset("")
}
