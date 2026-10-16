package orm

import (
	"fmt"
	"gin/config"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func openSqlsrv(cfg *config.Config, gormLog gormlogger.Interface) (*gorm.DB, error) {
	return gorm.Open(sqlserver.Open(getSqlsrvDsn(cfg)), &gorm.Config{
		NamingStrategy: configNaming(),
		Logger:         gormLog,
	})
}

func getSqlsrvDsn(cfg *config.Config) string {
	/*
	   官方推荐格式(最稳定)：
	   sqlserver://username:password@host:port?database=dbname
	   常见坑：
	   - password 有特殊字符需要 url.QueryEscape
	   - SQLServer 默认端口 1433
	*/

	return fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		cfg.Databases.Sqlsrv.Username,
		cfg.Databases.Sqlsrv.Password,
		cfg.Databases.Sqlsrv.Host,
		cfg.Databases.Sqlsrv.Port,
		cfg.Databases.Sqlsrv.Database,
	)
}
