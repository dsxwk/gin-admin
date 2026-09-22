package facade

import (
	"gin/pkg/serviceprovider/agent"
	"gin/pkg/serviceprovider/agent/providers"
)

// Agent 获取AI Agent实例
func Agent(name ...string) *AgentFacade {
	return (&AgentFacade{}).Provider(name...)
}

// AgentFacade AI Agent门面
type AgentFacade struct {
	*agent.Agent
}

// Provider 设置模型提供商
func (a *AgentFacade) Provider(name ...string) *AgentFacade {
	providerName, _, ok := AgentProvider(name...)
	if !ok {
		return nil
	}

	cfg := Config()
	if cfg == nil {
		return nil
	}

	providerCfg, exists := cfg.Agent.Providers[providerName]
	if !exists {
		return nil
	}

	provider := providers.NewOpenAICompat(providerName, providerCfg, cfg.Agent.MaxTokens, cfg.Agent.Temperature)
	facade := a
	if facade == nil {
		facade = &AgentFacade{}
	}
	facade.Agent = agent.New(provider, MCP().Tools())

	return facade
}

// AgentProvider 解析AI提供商信息,返回提供商名称和模型名
func AgentProvider(name ...string) (providerName, model string, ok bool) {
	cfg := Config()
	if cfg == nil || !cfg.Agent.Enabled {
		return "", "", false
	}

	providerName = cfg.Agent.Default
	if len(name) > 0 && name[0] != "" {
		providerName = name[0]
	}

	providerCfg, exists := cfg.Agent.Providers[providerName]
	if !exists {
		return "", "", false
	}

	return providerName, providerCfg.Model, true
}
