package provider

import (
	"gin/app/facade"
	"gin/pkg/foundation"
	"gin/pkg/orm"
	"go.uber.org/zap"
)

func init() {
	foundation.Register(&DbProvider{})
}

// DbProvider 数据库服务提供者
type DbProvider struct{}

// Name 服务提供者名称
func (p *DbProvider) Name() string {
	return "db"
}

// Register 注册服务到门面
func (p *DbProvider) Register(app foundation.App) {
	facade.Register("db", orm.Connection())
}

// Boot 启动服务-测试数据库连接
func (p *DbProvider) Boot(app foundation.App) {
	// 测试默认连接是否正常
	db := orm.Connection()
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err = sqlDB.Ping(); err == nil {
				facade.Log.Info("数据库连接成功", zap.String("driver", "mysql"))
			}
		}
	}
}

// Dependencies 依赖配置和日志服务
func (p *DbProvider) Dependencies() []string {
	return []string{"config", "log"}
}
