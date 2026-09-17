package provider

import (
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/http"
)

func init() {
	serviceprovider.Register(&HttpProvider{})
}

// HttpProvider HTTP客户端服务提供者
type HttpProvider struct{}

// Name 服务提供者名称
func (p *HttpProvider) Name() string {
	return serviceprovider.ServiceHTTP
}

// Register 注册服务到容器
func (p *HttpProvider) Register(app *container.Container) {
	app.Set(serviceprovider.ServiceHTTP, http.NewClient())
}

// Boot 启动服务
func (p *HttpProvider) Boot(app *container.Container) {
	http.GetClient()
	flag.Infof("HTTP客户端服务启动成功")
}

// Dependencies 依赖的服务
func (p *HttpProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceLog}
}
