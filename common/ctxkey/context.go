package ctxkey

import "context"

const (
	// DBKey 数据库连接键
	DBKey       string = "db_default" // 默认数据库连接
	DBMySQLKey  string = "db_mysql"   // MySQL连接
	DBPgsqlKey  string = "db_pgsql"   // PostgreSQL连接
	DBSqliteKey string = "db_sqlite"  // SQLite连接
	DBSqlsrvKey string = "db_sqlsrv"  // SQL Server连接

	// CacheKey 缓存键
	CacheKey       string = "cache_default" // 默认缓存
	CacheRedisKey  string = "cache_redis"   // Redis缓存
	CacheMemoryKey string = "cache_memory"  // 内存缓存
	CacheDiskKey   string = "cache_disk"    // 磁盘缓存

	// LogKey 日志键
	LogKey       string = "log"
	UserIdKey    string = "userId"
	ContainerKey string = "container"
	TraceIdKey   string = "traceId"
	IpKey        string = "ip"
	PathKey      string = "path"
	MethodKey    string = "method"
	ParamsKey    string = "params"
	MsKey        string = "ms"
	LangKey      string = "lang"
	StartTimeKey string = "startTime"
	DebuggerKey  string = "debugger"
)

// WithValue 将值注入到context
func WithValue(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

// GetValue 从context获取值
func GetValue(ctx context.Context, key string) interface{} {
	return ctx.Value(key)
}
