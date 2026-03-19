package connection

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
		conf.Pgsql.Host,
		conf.Pgsql.Username,
		conf.Pgsql.Password,
		conf.Pgsql.Database,
		conf.Pgsql.Port,
	)
}
