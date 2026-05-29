package provider

import (
	"gin/app/facade"
	"gin/common/flag"
	"gin/pkg/serviceprovider"
)

func init() {
	serviceprovider.Register(&RequestProvider{})
}

// RequestProvider 请求验证服务提供者
type RequestProvider struct{}

// Name 服务提供者名称
func (p *RequestProvider) Name() string {
	return "request"
}

// Register 注册服务到门面
func (p *RequestProvider) Register(app serviceprovider.App) {
	facade.Register("request", facade.Request[any]())
}

// Boot 启动服务
func (p *RequestProvider) Boot(app serviceprovider.App) {
	flag.Infof("请求验证服务启动成功")
}

// Dependencies 依赖服务
func (p *RequestProvider) Dependencies() []string {
	return []string{"log"} // 依赖日志服务
}
