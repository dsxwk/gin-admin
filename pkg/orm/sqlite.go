package orm

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSqlite() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(conf.Sqlite.Path), &gorm.Config{
		NamingStrategy: configNaming(),
		Logger:         gormLogger(),
	})
}
