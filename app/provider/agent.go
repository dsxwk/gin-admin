package provider

import (
	"gin/config"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	serviceagent "gin/pkg/serviceprovider/agent"
	"gin/pkg/serviceprovider/agent/providers"
	"gin/pkg/serviceprovider/mcp"
	"gin/pkg/terminal"
)

// AgentProvider AI Agent服务提供者
type AgentProvider struct {
	manager *serviceagent.Manager
}

// Name 服务提供者名称
func (p *AgentProvider) Name() string {
	return serviceprovider.ServiceAgent
}

// Register 注册服务到容器
func (p *AgentProvider) Register(app *container.Container) {
	cfg := app.Config()
	if cfg == nil || !cfg.Agent.Enabled {
		return
	}

	var tools []mcp.Tool
	if handler := app.MCP(); handler != nil {
		tools = handler.Tools()
	}

	p.manager = serviceagent.NewManager(cfg, tools, func(providerName string, providerConfig config.AgentProvider) serviceagent.Provider {
		return providers.NewOpenAICompat(
			providerName,
			providerConfig,
			cfg.Agent.MaxTokens,
			cfg.Agent.Temperature,
			cfg.Agent.RequestTimeout,
		)
	})
	app.SetAgent(p.manager)
}

// Boot 启动服务
func (p *AgentProvider) Boot(_ *container.Container) {
	if p.manager == nil {
		return
	}
	terminal.Infof("AI Agent服务启动成功")
}

// Dependencies 依赖服务
func (p *AgentProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig, serviceprovider.ServiceMCP}
}
