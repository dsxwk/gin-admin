package connection

import (
	"fmt"
	"gin/common/ctxkey"
	"gin/common/flag"
	"gin/config"
	"gin/pkg/debugger"
	l "gin/pkg/logger"
	"gin/pkg/message"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	conf        = config.NewConfig()
	dbInstances = make(map[string]*gorm.DB)
	dbLocks     sync.Map
)

type Db struct{}

// GetDB 初始化数据库(统一入口)
func (Db) GetDB() *gorm.DB {
	return getConnection(conf.Databases.DbConnection)
}

// Connection 连接数据库
func (Db) Connection(conn string) *gorm.DB {
	return getConnection(conn)
}

func getConnection(conn string) *gorm.DB {
	// 每个连接只初始化一次
	onceAny, _ := dbLocks.LoadOrStore(conn, &sync.Once{})
	once := onceAny.(*sync.Once)

	once.Do(func() {
		var (
			db  *gorm.DB
			err error
		)

		switch conn {

		case "mysql":
			db, err = openMysql()

		case "pgsql":
			db, err = openPgsql()

		case "sqlite":
			db, err = openSqlite()

		case "sqlsrv":
			db, err = openSqlsrv()

		default:
			color.Red(flag.Error+"  不支持的数据库类型: %s", conf.Databases.DbConnection)
			os.Exit(1)
		}

		if err != nil {
			color.Red(flag.Error+"  %s数据库连接失败: %v", conf.Databases.DbConnection, err)
			os.Exit(1)
		}

		// 配置连接池
		sqlDB, e := db.DB()
		if e != nil {
			err = e
			color.Red(flag.Error+"  %s数据库连接池配置失败: %v", conf.Databases.DbConnection, err)
			os.Exit(1)
		}

		// 设置连接池参数
		sqlDB.SetMaxIdleConns(20)                  // 空闲连接数
		sqlDB.SetMaxOpenConns(200)                 // 最大连接数
		sqlDB.SetConnMaxLifetime(1 * time.Hour)    // 连接最大生命周期
		sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 空闲时间超过30分钟自动关闭

		// 测试Ping
		if err = sqlDB.Ping(); e != nil {
			color.Red(flag.Error+"  %数据库连接ping失败: %v", conf.Databases.DbConnection, err)
			os.Exit(1)
		}

		// 注册gorm sql回调
		SqlCallback(db)

		dbInstances[conn] = db
	})

	return dbInstances[conn]
}

// configNaming 全局表名策略
func configNaming() schema.NamingStrategy {
	return schema.NamingStrategy{
		SingularTable: true, // 全局关闭表名复数化
	}
}

// gormLogger gorm日志配置
func gormLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到控制台
		logger.Config{
			SlowThreshold: conf.Mysql.SlowQueryDuration, // 慢sql阈值转Duration
			LogLevel:      logger.Info,
			Colorful:      true, // 彩色日志
			// IgnoreRecordNotFoundError: true, // 如果需要忽略 record not found
		},
	)
}

// SqlCallback sql回调
func SqlCallback(db *gorm.DB) {
	// 查询
	_ = db.Callback().Query().Before("gorm:query").Register("log:before_query", before)
	_ = db.Callback().Query().After("gorm:query").Register("log:after_query", after)

	// 创建
	_ = db.Callback().Create().Before("gorm:create").Register("log:before_create", before)
	_ = db.Callback().Create().After("gorm:create").Register("log:after_create", after)

	// 更新
	_ = db.Callback().Update().Before("gorm:update").Register("log:before_update", before)
	_ = db.Callback().Update().After("gorm:update").Register("log:after_update", after)

	// 删除
	_ = db.Callback().Delete().Before("gorm:delete").Register("log:before_delete", before)
	_ = db.Callback().Delete().After("gorm:delete").Register("log:after_delete", after)
}

func before(db *gorm.DB) {
	db.InstanceSet(startTimeKey, time.Now())
}

func after(db *gorm.DB) {
	ctx := db.Statement.Context
	if ctx == nil {
		return
	}
	start, ok := db.InstanceGet(startTimeKey)
	if !ok {
		return
	}

	// sql builder是否为空
	if db.Statement.SQL.Len() == 0 || db.Dialector == nil {
		return
	}

	var sql string
	// 安全处理Vars为空
	if db.Statement.Vars == nil || len(db.Statement.Vars) == 0 {
		sql = db.Statement.SQL.String()
	} else {
		sql = db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	}

	// 耗时
	cost := time.Since(start.(time.Time))
	costMs := float64(cost.Nanoseconds()) / 1e6 // 精确到小数

	// 慢查询警告
	if cost > conf.Mysql.SlowQueryDuration {
		l.NewLogger().Warn(
			"Slow Sql",
			zap.Float64("costMs", costMs),
			zap.String("sql", sql),
		)
	}

	traceId, ok := ctx.Value(ctxkey.TraceIdKey).(string)
	if !ok || traceId == "" {
		traceId = "unknown"
	}

	message.GetEventBus().Publish(debugger.TopicSql, debugger.SqlEvent{
		TraceId: traceId,
		Sql:     sql,
		Rows:    db.Statement.RowsAffected,
		Ms:      costMs,
	})
}

// getSql 替换Sql中的占位符"?"为实际值
func getSql(sql string, vars []interface{}) string {
	for _, v := range vars {
		// 将参数值格式化为字符串
		var (
			formattedValue string
		)
		switch value := v.(type) {
		case string:
			formattedValue = fmt.Sprintf("'%s'", value)
		case []byte:
			formattedValue = fmt.Sprintf("'%s'", string(value))
		case time.Time:
			formattedValue = fmt.Sprintf("'%s'", value.Format("2006-01-02 15:04:05"))
		case *time.Time:
			if value != nil {
				formattedValue = fmt.Sprintf("'%s'", value.Format("2006-01-02 15:04:05"))
			} else {
				formattedValue = "NULL"
			}
		case *gorm.DeletedAt:
			if value == nil {
				formattedValue = "NULL"
			}
		default:
			formattedValue = fmt.Sprintf("%v", value)
		}

		// 替换第一个"?"为实际值
		sql = strings.Replace(sql, "?", formattedValue, 1)
	}

	return sql
}
