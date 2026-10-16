package orm

import (

	oracle "github.com/godoes/gorm-oracle"
	"gorm.io/gorm"
	"fmt"
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
		connectString = fmt.Sprintf(
			"%s:%s/%s",
			conf.Databases.Oracle.Host,
			conf.Databases.Oracle.Port,
			conf.Databases.Oracle.ServiceName,
		)
	} else if conf.Databases.Oracle.Sid != "" {
		connectString = fmt.Sprintf(
			"%s:%s/%s?sid=true",
			conf.Databases.Oracle.Host,
			conf.Databases.Oracle.Port,
			conf.Databases.Oracle.Sid,
		)
	}

	return fmt.Sprintf(
		`user="%s" password="%s" connectString="%s"`,
		conf.Databases.Oracle.Username,
		conf.Databases.Oracle.Password,
		connectString,
	)
}
