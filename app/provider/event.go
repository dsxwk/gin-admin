package provider

import (
	"gin/app/facade"
	"gin/pkg/foundation"
)

func init() {
	foundation.Register(&EventProvider{})
}

// EventProvider 事件服务提供者
type EventProvider struct{}

func (p *EventProvider) Name() string {
	return "event"
}

func (p *EventProvider) Register(app foundation.App) {
	// 注册事件门面
	facade.Register("event", facade.Event)
}

func (p *EventProvider) Boot(app foundation.App) {
	facade.Log.Info("事件总线服务启动成功")
}

func (p *EventProvider) Dependencies() []string {
	return []string{"log"} // 依赖日志服务
}
