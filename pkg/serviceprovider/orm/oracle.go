package orm

import (
	"fmt"
	"gin/config"

	oracle "github.com/godoes/gorm-oracle"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func openOracle(cfg *config.Config, gormLog gormlogger.Interface) (*gorm.DB, error) {
	return gorm.Open(oracle.Open(getOracleDsn(cfg)), &gorm.Config{
		NamingStrategy: configNaming(),
		Logger:         gormLog,
	})
}

func getOracleDsn(cfg *config.Config) string {
	connectString := ""
	if cfg.Databases.Oracle.ServiceName != "" {
		connectString = fmt.Sprintf(
			"%s:%s/%s",
			cfg.Databases.Oracle.Host,
			cfg.Databases.Oracle.Port,
			cfg.Databases.Oracle.ServiceName,
		)
	} else if cfg.Databases.Oracle.Sid != "" {
		connectString = fmt.Sprintf(
			"%s:%s/%s?sid=true",
			cfg.Databases.Oracle.Host,
			cfg.Databases.Oracle.Port,
			cfg.Databases.Oracle.Sid,
		)
	}

	return fmt.Sprintf(
		`user="%s" password="%s" connectString="%s"`,
		cfg.Databases.Oracle.Username,
		cfg.Databases.Oracle.Password,
		connectString,
	)
}
