package provider

import (
	"gin/pkg/serviceprovider"
)

// All 服务提供者列表
func All() []serviceprovider.ServiceProvider {
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
