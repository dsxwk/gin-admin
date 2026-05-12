package provider

import (
	"gin/app/facade"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/foundation"
	"gin/pkg/provider/orm"
	"gorm.io/gorm"
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
	cfg := facade.Config()
	facade.Register[*gorm.DB]("db", orm.Connection(cfg.Databases.Default, cfg))
}

// Boot 启动服务-测试数据库连接
func (p *DbProvider) Boot(app foundation.App) {
	cfg := facade.Config()
	// 测试默认连接是否正常
	db := orm.Connection(cfg.Databases.Default, cfg)
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err = sqlDB.Ping(); err == nil {
				flag.Infof(pkg.Sprintf("%s数据库连接成功", cfg.Databases.Default))
			}
		}
	}
}

// Dependencies 依赖配置和日志服务
func (p *DbProvider) Dependencies() []string {
	return []string{"config", "log"}
}
