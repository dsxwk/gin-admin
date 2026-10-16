package provider

import (
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/orm"
)

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
	p.manager = orm.NewManager(app.Config(), app.Log())
	app.SetDB(p.manager)
}

// Boot 启动服务-测试数据库连接
func (p *DbProvider) Boot(app *container.Container) {
	driver := app.Config().Databases.Driver
	if _, err := p.manager.Connect(driver); err != nil {
		flag.Errorf("%v", err)
		return
	}
	flag.Infof("%s数据库连接成功", driver)
}

// Dependencies 依赖配置和日志服务
func (p *DbProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig, serviceprovider.ServiceLog}
}
