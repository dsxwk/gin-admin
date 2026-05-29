package provider

import (
	"gin/app/facade"
	"gin/common/flag"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/lang"
)

func init() {
	serviceprovider.Register(&LangProvider{})
}

// LangProvider 翻译服务提供者
type LangProvider struct{}

// Name 服务提供者名称
func (p *LangProvider) Name() string {
	return "lang"
}

// Register 注册服务到门面
func (p *LangProvider) Register(app serviceprovider.App) {
	// 翻译服务在Boot时加载,这里只做占位
	facade.Register("lang", facade.Lang())
}

// Boot 启动服务
func (p *LangProvider) Boot(app serviceprovider.App) {
	// 加载翻译文件(会从facade.Config()获取配置)
	lang.LoadLang(facade.Config(), facade.Log())
	flag.Infof("翻译服务启动成功")
}

// Dependencies 依赖服务
func (p *LangProvider) Dependencies() []string {
	return []string{"config", "log"} // 依赖配置和日志
}
