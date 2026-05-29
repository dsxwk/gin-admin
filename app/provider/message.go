package provider

import (
	"gin/app/facade"
	"gin/common/flag"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/message"
)

func init() {
	serviceprovider.Register(&MessageProvider{})
}

// MessageProvider 消息事件服务提供者
type MessageProvider struct{}

func (p *MessageProvider) Name() string {
	return "message"
}

func (p *MessageProvider) Register(app serviceprovider.App) {
	// 注册事件门面
	facade.Register[*message.Event]("message", message.NewEvent())
}

func (p *MessageProvider) Boot(app serviceprovider.App) {
	flag.Infof("消息事件服务启动成功")
}

func (p *MessageProvider) Dependencies() []string {
	return []string{"log"}
}
