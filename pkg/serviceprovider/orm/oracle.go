package orm

import (
	"gin/pkg"
	oracle "github.com/godoes/gorm-oracle"
	"gorm.io/gorm"
)

func openOracle() (*gorm.DB, error) {
	return gorm.Open(oracle.Open(getOracleDsn()), &gorm.Config{
		NamingStrategy: configNaming(),
		Logger:         gormLogger(),
	})
}

func getOracleDsn() string {
	connectString := ""
	if conf.Databases.Oracle.ServiceName != "" {
		connectString = pkg.Sprintf(
			"%s:%s/%s",
			conf.Databases.Oracle.Host,
			conf.Databases.Oracle.Port,
			conf.Databases.Oracle.ServiceName,
		)
	} else if conf.Databases.Oracle.Sid != "" {
		connectString = pkg.Sprintf(
			"%s:%s/%s?sid=true",
			conf.Databases.Oracle.Host,
			conf.Databases.Oracle.Port,
			conf.Databases.Oracle.Sid,
		)
	}

	return pkg.Sprintf(
		`user="%s" password="%s" connectString="%s"`,
		conf.Databases.Oracle.Username,
		conf.Databases.Oracle.Password,
		connectString,
	)
}
