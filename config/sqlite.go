package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSqlite() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(Conf.Sqlite.Path), &gorm.Config{
		NamingStrategy: configNaming(),
		Logger:         gormLogger(),
	})
}
