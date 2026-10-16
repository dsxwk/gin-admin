package orm

import (
	"fmt"
	"gin/config"
	"gin/pkg/ctxkey"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/eventbus"
	servicelogger "gin/pkg/serviceprovider/logger"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	startTimeKey           = "gorm_start_time"
	defaultMaxIdleConns    = 20
	defaultMaxOpenConns    = 200
	defaultConnMaxLifetime = 3 * time.Minute
	defaultConnMaxIdleTime = 1 * time.Minute
)

// openConnection 根据驱动创建数据库连接
func openConnection(driver string, cfg *config.Config, gormLog gormlogger.Interface) (*gorm.DB, error) {
	if cfg == nil {
		return nil, fmt.Errorf("数据库配置未初始化")
	}

	switch driver {
	case "mysql":
		return openMysql(cfg, gormLog)
	case "pgsql":
		return openPgsql(cfg, gormLog)
	case "sqlite":
		return openSqlite(cfg, gormLog)
	case "sqlsrv":
		return openSqlsrv(cfg, gormLog)
	case "oracle":
		return openOracle(cfg, gormLog)
	default:
		return nil, fmt.Errorf("不支持的数据库: %s", driver)
	}
}

// connectionPoolSettings 获取连接池参数,未配置时使用默认值
func connectionPoolSettings(cfg *config.Config) (int, int, time.Duration, time.Duration) {
	if cfg == nil {
		return defaultMaxIdleConns, defaultMaxOpenConns, defaultConnMaxLifetime, defaultConnMaxIdleTime
	}

	pool := cfg.Databases.Pool
	maxIdleConns := pool.MaxIdleConns
	if maxIdleConns <= 0 {
		maxIdleConns = defaultMaxIdleConns
	}

	maxOpenConns := pool.MaxOpenConns
	if maxOpenConns <= 0 {
		maxOpenConns = defaultMaxOpenConns
	}

	connMaxLifetime := pool.ConnMaxLifetime
	if connMaxLifetime <= 0 {
		connMaxLifetime = defaultConnMaxLifetime
	}

	connMaxIdleTime := pool.ConnMaxIdleTime
	if connMaxIdleTime <= 0 {
		connMaxIdleTime = defaultConnMaxIdleTime
	}

	return maxIdleConns, maxOpenConns, connMaxLifetime, connMaxIdleTime
}

// configNaming 全局表名策略
func configNaming() schema.NamingStrategy {
	return schema.NamingStrategy{
		SingularTable: true,
	}
}

// newGormLogger 创建GORM日志配置
func newGormLogger(cfg *config.Config) gormlogger.Interface {
	slowThreshold := 3 * time.Second
	if cfg != nil && cfg.Databases.SlowQueryDuration > 0 {
		slowThreshold = cfg.Databases.SlowQueryDuration
	}

	return gormlogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormlogger.Config{
			SlowThreshold: slowThreshold,
			LogLevel:      gormlogger.Info,
			Colorful:      true,
		},
	)
}

// sqlCallback 注册SQL回调
func sqlCallback(db *gorm.DB, cfg *config.Config, log *servicelogger.Logger) {
	afterCallback := func(db *gorm.DB) {
		after(db, cfg, log)
	}

	_ = db.Callback().Query().Before("gorm:query").Register("log:before_query", before)
	_ = db.Callback().Query().After("gorm:query").Register("log:after_query", afterCallback)

	_ = db.Callback().Create().Before("gorm:create").Register("log:before_create", before)
	_ = db.Callback().Create().After("gorm:create").Register("log:after_create", afterCallback)

	_ = db.Callback().Update().Before("gorm:update").Register("log:before_update", before)
	_ = db.Callback().Update().After("gorm:update").Register("log:after_update", afterCallback)

	_ = db.Callback().Delete().Before("gorm:delete").Register("log:before_delete", before)
	_ = db.Callback().Delete().After("gorm:delete").Register("log:after_delete", afterCallback)
}

// before SQL执行前记录开始时间
func before(db *gorm.DB) {
	db.InstanceSet(startTimeKey, time.Now())
}

// after SQL执行后记录慢查询和调试事件
func after(db *gorm.DB, cfg *config.Config, log *servicelogger.Logger) {
	ctx := db.Statement.Context
	if ctx == nil {
		return
	}

	start, ok := db.InstanceGet(startTimeKey)
	if !ok {
		return
	}
	if db.Statement.SQL.Len() == 0 || db.Dialector == nil {
		return
	}

	var sql string
	if db.Statement.Vars == nil || len(db.Statement.Vars) == 0 {
		sql = db.Statement.SQL.String()
	} else {
		sql = db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	}

	cost := time.Since(start.(time.Time))
	costMs := float64(cost.Nanoseconds()) / 1e6

	if cfg != nil && log != nil && cost > cfg.Databases.SlowQueryDuration {
		log.Warn(
			"Slow Sql",
			zap.Float64("costMs", costMs),
			zap.String("sql", sql),
		)
	}

	eventbus.Default().Publish(ctx, debugger.TopicSQL, debugger.SQLEvent{
		TraceID: ctxkey.TraceID(ctx),
		SQL:     sql,
		Rows:    db.Statement.RowsAffected,
		Ms:      costMs,
	})
}
