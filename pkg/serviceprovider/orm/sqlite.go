package orm

import (
	"gin/config"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	// 纯Go的SQLite驱动,不需要CGO
	"github.com/glebarez/sqlite"
)

func openSqlite(cfg *config.Config, gormLog gormlogger.Interface) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(cfg.Databases.Sqlite.Path), &gorm.Config{
		NamingStrategy: configNaming(),
		Logger:         gormLog,
	})
}
