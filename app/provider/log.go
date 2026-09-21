package provider

import (
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/logger"
)

// LogProvider 日志服务提供者
type LogProvider struct{}

// Name 服务提供者名称
func (p *LogProvider) Name() string {
	return serviceprovider.ServiceLog
}

// Register 注册日志服务到容器
func (p *LogProvider) Register(app *container.Container) {
	app.SetLog(logger.NewLogger(app.Config()))
}

// Boot 启动服务
func (p *LogProvider) Boot(app *container.Container) {
	flag.Infof("日志服务启动成功")
}

// Dependencies 依赖服务
func (p *LogProvider) Dependencies() []string {
	return []string{"config"}
}
