package provider

import (
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/request"
)

func init() {
	serviceprovider.Register(&RequestProvider{})
}

// RequestProvider 请求验证服务提供者
type RequestProvider struct{}

// Name 服务提供者名称
func (p *RequestProvider) Name() string {
	return serviceprovider.ServiceRequest
}

// Register 注册服务到容器
func (p *RequestProvider) Register(app *container.Container) {
	app.Set(serviceprovider.ServiceRequest, request.NewClient())
}

// Boot 启动服务
func (p *RequestProvider) Boot(app *container.Container) {
	flag.Infof("请求验证服务启动成功")
}

// Dependencies 依赖服务
func (p *RequestProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceLog}
}
