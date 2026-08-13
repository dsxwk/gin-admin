package provider

import (
	"gin/app/facade"
	"gin/common/flag"
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
	return "mcp"
}

// Register 注册服务到门面
func (p *McpProvider) Register(app serviceprovider.App) {
	cfg := facade.Config()
	if cfg == nil || !cfg.Mcp.Enabled {
		return
	}

	handler := mcp.NewHandler(mcp.ServerInfo{
		Name:    cfg.App.Name,
		Version: cfg.App.CliVersion,
	}, cfg)

	facade.Register("mcp", handler)
	flag.Infof("MCP服务注册成功,路径: %s", cfg.Mcp.Path)
}

// Boot 启动服务
func (p *McpProvider) Boot(app serviceprovider.App) {
	// 路由在LoadRouters中注册
}

// Dependencies 依赖配置服务
func (p *McpProvider) Dependencies() []string {
	return []string{"config"}
}
