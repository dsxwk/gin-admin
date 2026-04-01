package provider

import (
	"gin/app/facade"
	"gin/config"
	"gin/pkg/foundation"
)

func init() {
	foundation.Register(&ConfigProvider{})
}

// ConfigProvider 配置服务提供者
type ConfigProvider struct{}

// Name 服务提供者名称
func (p *ConfigProvider) Name() string {
	return "config"
}

// Register 注册服务到门面
func (p *ConfigProvider) Register(app foundation.App) {
	// 获取配置实例
	cfg := config.NewConfig()
	// 注册到门面
	facade.Register("config", cfg)
}

// Boot 启动服务(配置服务无需额外启动逻辑)
func (p *ConfigProvider) Boot(app foundation.App) {}

// Dependencies 依赖服务
func (p *ConfigProvider) Dependencies() []string {
	return []string{}
}
