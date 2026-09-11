package provider

import (
	"gin/app/facade"
	"gin/common/flag"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/eventbus"
)

func init() {
	serviceprovider.Register(&EventProvider{})
}

// EventProvider 事件服务提供者
type EventProvider struct{}

func (p *EventProvider) Name() string {
	return "event"
}

func (p *EventProvider) Register(app serviceprovider.App) {
	// 注册事件门面
	facade.Register[*eventbus.Bus]("event", eventbus.NewBus())
}

func (p *EventProvider) Boot(app serviceprovider.App) {
	flag.Infof("事件服务启动成功")
}

func (p *EventProvider) Dependencies() []string {
	return []string{"log"} // 依赖日志服务
}
