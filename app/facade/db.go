package facade

import (
	"gin/pkg/orm"
	"gorm.io/gorm"
)

// DB 数据库门面-数据库访问统一入口
var DB = &dbFacade{}

type dbFacade struct{}

// Connection 获取指定名称的数据库连接
func (d *dbFacade) Connection(connName ...string) *gorm.DB {
	return orm.Connection(connName...)
}
