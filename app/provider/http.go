package provider

import (
	"gin/app/facade"
	"gin/pkg/foundation"
	"gin/pkg/provider/http"
)

func init() {
	foundation.Register(&HttpProvider{})
}

// HttpProvider HTTP客户端服务提供者
type HttpProvider struct{}

// Name 服务提供者名称
func (p *HttpProvider) Name() string {
	return "http"
}

// Register 注册服务到门面
func (p *HttpProvider) Register(app foundation.App) {
	facade.Register[*http.Client[any]]("http", http.NewClient[any]())
}

// Boot 启动服务
func (p *HttpProvider) Boot(app foundation.App) {
	http.InitClient()
	facade.Log().Info("HTTP客户端服务启动成功")
}

// Dependencies 依赖的服务
func (p *HttpProvider) Dependencies() []string {
	return []string{"log"}
}
