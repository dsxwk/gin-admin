package provider

import (
	"gin/app/listener"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/eventbus"
	"gin/pkg/terminal"
)

// EventProvider 事件服务提供者
type EventProvider struct{}

// Name 服务提供者名称
func (p *EventProvider) Name() string {
	return serviceprovider.ServiceEvent
}

// Register 创建事件总线并注册业务事件表
func (p *EventProvider) Register(app *container.Container) {
	registry := eventbus.NewRegistry(eventbus.Default())
	app.SetEvent(registry)
	eventbus.Register(registry, listener.All())
}

// Boot 启动事件服务
func (p *EventProvider) Boot(_ *container.Container) {
	terminal.Infof("事件服务启动成功")
}

// Dependencies 依赖服务
func (p *EventProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceLog}
}
