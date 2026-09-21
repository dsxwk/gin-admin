package provider

import (
	"gin/pkg/serviceprovider"
)

// Providers 获取服务提供者列表
func Providers() []serviceprovider.ServiceProvider {
	return []serviceprovider.ServiceProvider{
		&ConfigProvider{},
		&LogProvider{},
		&DbProvider{},
		&CacheProvider{},
		&EventProvider{},
		&DebuggerProvider{},
		&RequestProvider{},
		&HttpProvider{},
		&EsProvider{},
		&GrpcProvider{},
		&QueueProvider{},
		&JobProvider{},
		&LangProvider{},
		&McpProvider{},
		&RateLimitProvider{},
	}
}
