package orm

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
		conf.Databases.Sqlsrv.Username,
		conf.Databases.Sqlsrv.Password,
		conf.Databases.Sqlsrv.Host,
		conf.Databases.Sqlsrv.Port,
		conf.Databases.Sqlsrv.Database,
	)
}
