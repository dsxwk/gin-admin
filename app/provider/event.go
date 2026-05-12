package provider

import (
	"gin/app/facade"
	"gin/common/flag"
	"gin/pkg/foundation"
	"gin/pkg/provider/message"
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
	facade.Register[*message.Event]("event", facade.Message())
}

func (p *EventProvider) Boot(app foundation.App) {
	flag.Infof("事件服务启动成功")
}

func (p *EventProvider) Dependencies() []string {
	return []string{"message", "log"} // 依赖日志服务
}
