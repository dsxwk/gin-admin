package provider

import (
	"gin/app/facade"
	"gin/pkg/foundation"
	"gin/pkg/lang"
)

func init() {
	foundation.Register(&LangProvider{})
}

// LangProvider 翻译服务提供者
type LangProvider struct{}

// Name 服务提供者名称
func (p *LangProvider) Name() string {
	return "lang"
}

// Register 注册服务到门面
func (p *LangProvider) Register(app foundation.App) {
	// 翻译服务在Boot时加载,这里只做占位
	facade.Register("lang", facade.Lang)
}

// Boot 启动服务
func (p *LangProvider) Boot(app foundation.App) {
	// 加载翻译文件(会从facade.Config获取配置)
	lang.LoadLang(facade.Config.Get(), facade.Log.Logger())
	facade.Log.Info("翻译服务启动成功")
}

// Dependencies 依赖服务
func (p *LangProvider) Dependencies() []string {
	return []string{"config", "log"} // 依赖配置和日志
}
