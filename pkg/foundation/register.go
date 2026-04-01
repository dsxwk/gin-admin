package foundation

import (
	"fmt"
	"gin/common/flag"
	"os"
	"sync"
)

// Registry 服务提供者注册表
type Registry struct {
	mu        sync.RWMutex
	providers []ServiceProvider
	names     map[string]bool // 记录已注册的提供者名称,防止重复注册
}

var (
	globalRegistry = &Registry{
		providers: make([]ServiceProvider, 0),
		names:     make(map[string]bool),
	}
)

// Register 注册服务提供者(由各个提供者的init函数调用,实现自动注册)
// 使用示例: foundation.Register(&ConfigProvider{})
func Register(provider ServiceProvider) {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	name := provider.Name()
	// 防止重复注册
	if globalRegistry.names[name] {
		flag.Errorf("Provider %s already registered", name)
		os.Exit(1)
	}

	globalRegistry.providers = append(globalRegistry.providers, provider)
	globalRegistry.names[name] = true
}

// GetProviders 获取所有已注册的服务提供者
func GetProviders() []ServiceProvider {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	// 返回副本,避免外部修改
	providers := make([]ServiceProvider, len(globalRegistry.providers))
	copy(providers, globalRegistry.providers)
	return providers
}

// getProviderName 获取服务提供者名称
func getProviderName(provider ServiceProvider) string {
	if named, ok := provider.(interface{ Name() string }); ok {
		return named.Name()
	}
	return fmt.Sprintf("%T", provider)
}
