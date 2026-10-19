package agent

import (
	"fmt"

	"gin/config"
	"gin/pkg/serviceprovider/mcp"
)

// ProviderFactory AI提供商工厂
type ProviderFactory func(providerName string, providerConfig config.AgentProvider) Provider

// Manager AI Agent管理器
type Manager struct {
	config   *config.Config
	tools    []mcp.Tool
	provider ProviderFactory
}

// NewManager 创建AI Agent管理器
func NewManager(cfg *config.Config, tools []mcp.Tool, provider ProviderFactory) *Manager {
	return &Manager{
		config:   cfg,
		tools:    append([]mcp.Tool(nil), tools...),
		provider: provider,
	}
}

// Provider 解析AI提供商信息
func (m *Manager) Provider(name ...string) (providerName, model string, ok bool) {
	if m == nil || m.config == nil || !m.config.Agent.Enabled {
		return "", "", false
	}

	providerName = m.config.Agent.Default
	if len(name) > 0 && name[0] != "" {
		providerName = name[0]
	}

	providerConfig, exists := m.config.Agent.Providers[providerName]
	if !exists {
		return "", "", false
	}
	return providerName, providerConfig.Model, true
}

// Tools 获取可用工具
func (m *Manager) Tools() []mcp.Tool {
	if m == nil {
		return nil
	}
	return append([]mcp.Tool(nil), m.tools...)
}

// Create 创建AI Agent实例
func (m *Manager) Create(name ...string) (*Agent, string, string, error) {
	if m == nil || m.config == nil || !m.config.Agent.Enabled {
		return nil, "", "", fmt.Errorf("AI Agent未启用")
	}

	providerName, model, ok := m.Provider(name...)
	if !ok {
		return nil, "", "", fmt.Errorf("AI Agent提供商不存在")
	}
	if m.provider == nil {
		return nil, "", "", fmt.Errorf("AI Agent提供商工厂未初始化")
	}

	providerConfig := m.config.Agent.Providers[providerName]
	provider := m.provider(providerName, providerConfig)
	if provider == nil {
		return nil, "", "", fmt.Errorf("AI Agent提供商创建失败")
	}

	return New(provider, m.Tools()), providerName, model, nil
}
