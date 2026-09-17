package provider

import (
	"gin/config"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
)

func init() {
	serviceprovider.Register(&ConfigProvider{})
}

// ConfigProvider 配置服务提供者
type ConfigProvider struct{}

// Name 服务提供者名称
func (p *ConfigProvider) Name() string {
	return serviceprovider.ServiceConfig
}

// Register 注册服务到容器
func (p *ConfigProvider) Register(app *container.Container) {
	cfg := config.NewConfig()
	app.Set(serviceprovider.ServiceConfig, cfg)

	config.OnConfigUpdated = func(updated *config.Config) {
		app.Set(serviceprovider.ServiceConfig, updated)
	}
}

// Boot 启动服务(配置服务无需额外启动逻辑)
func (p *ConfigProvider) Boot(app *container.Container) {}

// Dependencies 依赖服务
func (p *ConfigProvider) Dependencies() []string {
	return []string{}
}
