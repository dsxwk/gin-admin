package config

import (
	"gin/pkg"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func openSqlsrv() (*gorm.DB, error) {
	return gorm.Open(sqlserver.Open(getSqlsrvDsn()), &gorm.Config{
		NamingStrategy: configNaming(),
		Logger:         gormLogger(),
	})
}

func getSqlsrvDsn() string {
	/*
	   官方推荐格式(最稳定)：
	   sqlserver://username:password@host:port?database=dbname
	   常见坑：
	   - password 有特殊字符需要 url.QueryEscape
	   - SQLServer 默认端口 1433
	*/

	return pkg.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		Conf.Sqlsrv.Username,
		Conf.Sqlsrv.Password,
		Conf.Sqlsrv.Host,
		Conf.Sqlsrv.Port,
		Conf.Sqlsrv.Database,
	)
}
