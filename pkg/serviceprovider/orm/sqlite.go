package orm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	// 纯Go的SQLite驱动,不需要CGO
	"github.com/glebarez/sqlite"
)

func openSqlite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(conf.Databases.Sqlite.Path), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	return db, err
}
