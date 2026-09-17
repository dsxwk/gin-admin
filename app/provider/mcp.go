package provider

import (
	"gin/common/flag"
	"gin/config"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/mcp"
)

func init() {
	serviceprovider.Register(&McpProvider{})
}

// McpProvider MCP服务提供者
type McpProvider struct{}

// Name 服务提供者名称
func (p *McpProvider) Name() string {
	return serviceprovider.ServiceMCP
}

// Register 注册服务到容器
func (p *McpProvider) Register(app *container.Container) {
	cfg := app.Get[*config.Config](serviceprovider.ServiceConfig)
	if cfg == nil || !cfg.Mcp.Enabled {
		return
	}

	handler := mcp.NewHandler(mcp.ServerInfo{
		Name:    cfg.App.Name,
		Version: cfg.App.CliVersion,
	}, cfg)

	app.Set(serviceprovider.ServiceMCP, handler)
	flag.Infof("MCP服务注册成功,路径: %s", cfg.Mcp.Path)
}

// Boot 启动服务
func (p *McpProvider) Boot(app *container.Container) {
	// 路由在LoadRouters中注册
}

// Dependencies 依赖配置服务
func (p *McpProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig}
}
