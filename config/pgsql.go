package config

import (
	"gin/pkg"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openPgsql() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(getPgsqlDsn()), &gorm.Config{
		NamingStrategy: configNaming(),
		Logger:         gormLogger(),
	})
}

func getPgsqlDsn() string {
	return pkg.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		Conf.Pgsql.Host,
		Conf.Pgsql.Username,
		Conf.Pgsql.Password,
		Conf.Pgsql.Database,
		Conf.Pgsql.Port,
	)
}
