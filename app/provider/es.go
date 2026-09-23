package provider

import (
	"context"
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	esclient "gin/pkg/serviceprovider/es"
)

// EsProvider ES服务提供者
type EsProvider struct {
	client *esclient.Client
}

// Name 服务提供者名称
func (p *EsProvider) Name() string {
	return serviceprovider.ServiceES
}

// Register 注册服务到容器
func (p *EsProvider) Register(app *container.Container) {
	cfg := app.Config()
	if cfg == nil || !cfg.Es.Enabled {
		return
	}
	if len(cfg.Es.Addresses) == 0 {
		flag.Errorf("ES地址未配置")
		return
	}

	p.client = esclient.NewClient(cfg.Es)
	app.SetES(p.client)
}

// Boot 启动服务
func (p *EsProvider) Boot(_ *container.Container) {
	if p.client == nil {
		return
	}
	if err := p.client.Ping(context.Background()); err != nil {
		flag.Errorf("ES服务连接失败: %v", err)
		return
	}
	flag.Infof("ES服务启动成功")
}

// Dependencies 依赖服务
func (p *EsProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig, serviceprovider.ServiceLog, serviceprovider.ServiceHTTP}
}
