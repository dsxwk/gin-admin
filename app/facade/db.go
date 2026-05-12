package facade

import (
	"gin/pkg/provider/orm"
	"gorm.io/gorm"
)

// DB 数据库门面-数据库访问统一入口
func DB(conn ...string) *gorm.DB {
	name := "mysql"
	if len(conn) > 0 && conn[0] != "" {
		name = conn[0]
	}

	db := Get[*gorm.DB](name)
	if db != nil {
		return db
	}
	// orm.SetConfig(facade.Config())
	return orm.Connection(conn...)
}
