package provider

import (
	"gin/common/flag"
	"gin/config"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/orm"
)

func init() {
	serviceprovider.Register(&DbProvider{})
}

// DbProvider 数据库服务提供者
type DbProvider struct {
	manager *orm.Manager
}

// Name 服务提供者名称
func (p *DbProvider) Name() string {
	return serviceprovider.ServiceDB
}

// Register 注册服务到容器
func (p *DbProvider) Register(app *container.Container) {
	cfg := app.Get[*config.Config](serviceprovider.ServiceConfig)
	p.manager = orm.NewManager(cfg)
	app.Set(serviceprovider.ServiceDB, p.manager)
}

// Boot 启动服务-测试数据库连接
func (p *DbProvider) Boot(app *container.Container) {
	cfg := app.Get[*config.Config](serviceprovider.ServiceConfig)
	p.manager.Connection(cfg.Databases.Driver)
	flag.Infof("%s数据库连接成功", cfg.Databases.Driver)
}

// Dependencies 依赖配置和日志服务
func (p *DbProvider) Dependencies() []string {
	return []string{"config", "log"}
}
