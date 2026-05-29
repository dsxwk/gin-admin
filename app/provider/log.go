package provider

import (
	"gin/app/facade"
	"gin/common/flag"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/logger"
)

func init() {
	serviceprovider.Register(&LogProvider{})
}

// LogProvider 日志服务提供者
type LogProvider struct{}

// Name 服务提供者名称
func (p *LogProvider) Name() string {
	return "log"
}

// Register 注册日志服务到门面
func (p *LogProvider) Register(app serviceprovider.App) {
	// 注册到门面
	facade.Register[*logger.Logger]("log", logger.NewLogger(facade.Config()))
}

// Boot 启动服务
func (p *LogProvider) Boot(app serviceprovider.App) {
	flag.Infof("日志服务启动成功")
}

// Dependencies 依赖服务
func (p *LogProvider) Dependencies() []string {
	return []string{"config"}
}
