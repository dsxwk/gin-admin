package orm

import (
	"fmt"
	"gin/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func openPgsql(cfg *config.Config, gormLog gormlogger.Interface) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(getPgsqlDsn(cfg)), &gorm.Config{
		NamingStrategy: configNaming(),
		Logger:         gormLog,
	})
}

func getPgsqlDsn(cfg *config.Config) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Databases.Pgsql.Host,
		cfg.Databases.Pgsql.Username,
		cfg.Databases.Pgsql.Password,
		cfg.Databases.Pgsql.Database,
		cfg.Databases.Pgsql.Port,
		cfg.Databases.Pgsql.SSLMode,
	)
}
