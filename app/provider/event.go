package provider

import (
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/eventbus"
)

func init() {
	serviceprovider.Register(&EventProvider{})
}

// EventProvider 事件服务提供者
type EventProvider struct{}

// Name 服务提供者名称
func (p *EventProvider) Name() string {
	return serviceprovider.ServiceEvent
}

// Register 创建事件总线并注册业务事件表
func (p *EventProvider) Register(app *container.Container) {
	registry := eventbus.DefaultRegistry()
	bus := registry.Bus()

	app.Set(serviceprovider.ServiceEvent, registry)
	app.Set("eventbus", bus)
}

// Boot 启动事件服务
func (p *EventProvider) Boot(app *container.Container) {
	flag.Infof("事件服务启动成功")
}

// Dependencies 依赖服务
func (p *EventProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceLog}
}
