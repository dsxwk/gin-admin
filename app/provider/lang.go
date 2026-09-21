package provider

import (
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/lang"
)

// LangProvider 翻译服务提供者
type LangProvider struct {
	service *lang.Service
}

// Name 服务提供者名称
func (p *LangProvider) Name() string {
	return serviceprovider.ServiceLang
}

// Register 注册服务到容器
func (p *LangProvider) Register(app *container.Container) {
	p.service = lang.New()
	app.SetLang(p.service)
}

// Boot 启动服务
func (p *LangProvider) Boot(app *container.Container) {
	p.service.Load(app.Config(), app.Log())
	flag.Infof("翻译服务启动成功")
}

// Dependencies 依赖服务
func (p *LangProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig, serviceprovider.ServiceLog}
}
