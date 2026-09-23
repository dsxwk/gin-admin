package provider

import (
	appmcp "gin/app/mcp"
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	servicemcp "gin/pkg/serviceprovider/mcp"
)

// McpProvider MCP服务提供者
type McpProvider struct{}

// Name 服务提供者名称
func (p *McpProvider) Name() string {
	return serviceprovider.ServiceMCP
}

// Register 注册服务到容器
func (p *McpProvider) Register(app *container.Container) {
	cfg := app.Config()
	if cfg == nil || !cfg.Mcp.Enabled {
		return
	}

	handler, err := servicemcp.NewHandler(servicemcp.ServerInfo{
		Name:    cfg.App.Name,
		Version: cfg.App.CliVersion,
	}, cfg, appmcp.All())
	if err != nil {
		flag.Errorf("MCP服务注册失败: %v", err)
		return
	}

	app.SetMCP(handler)
	flag.Infof("MCP服务注册成功,路径: %s", cfg.Mcp.Path)
}

// Boot 启动服务
func (p *McpProvider) Boot(_ *container.Container) {
	// 路由在LoadRouters中注册
}

// Dependencies 依赖配置服务
func (p *McpProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig}
}
