package provider

import (
	"gin/app/facade"
	"gin/pkg/foundation"
)

func init() {
	foundation.Register(&MessageProvider{})
}

// MessageProvider 消息事件服务提供者
type MessageProvider struct{}

func (p *MessageProvider) Name() string {
	return "message"
}

func (p *MessageProvider) Register(app foundation.App) {
	// 注册事件门面
	facade.Register("message", facade.Message)
}

func (p *MessageProvider) Boot(app foundation.App) {
	facade.Log.Info("消息事件总线服务启动成功")
}

func (p *MessageProvider) Dependencies() []string {
	return []string{"log"}
}
