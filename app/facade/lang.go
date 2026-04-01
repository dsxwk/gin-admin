package facade

import (
	"context"
	"gin/pkg/lang"
)

// Lang 翻译门面-翻译统一入口
// 使用示例:
//
//	msg := facade.Lang.T(ctx, "login.success", map[string]interface{}{"name": "admin"})
//	msg := facade.Lang.T(ctx, "login.accountErr", nil)
var Lang = &langFacade{}

type langFacade struct{}

// T 翻译
// 参数:
//   - ctx: 上下文,用于获取语言设置
//   - messageID: 消息ID,如 "login.success"
//   - data: 模板数据,如 map[string]interface{}{"name": "admin"}
//
// 返回: 翻译后的字符串
func (l *langFacade) T(ctx context.Context, messageID string, data map[string]interface{}) string {
	return lang.T(ctx, messageID, data)
}

// GetLocalizer 获取指定语言
func (l *langFacade) GetLocalizer(langCode string) interface{} {
	return lang.GetLocalizer(langCode)
}

// GetBundle 获取翻译包
func (l *langFacade) GetBundle() interface{} {
	return lang.GetBundle()
}
