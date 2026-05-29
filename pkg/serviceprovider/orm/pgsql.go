package orm

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
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		conf.Databases.Pgsql.Host,
		conf.Databases.Pgsql.Username,
		conf.Databases.Pgsql.Password,
		conf.Databases.Pgsql.Database,
		conf.Databases.Pgsql.Port,
		conf.Databases.Pgsql.SSLMode,
	)
}
