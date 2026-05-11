package facade

import (
	"context"
	"gin/pkg/provider/lang"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Lang 翻译门面
// 使用示例:
//
//	msg := facade.Lang().T(ctx, "welcome", map[string]interface{}{"name": "John"})
//	localizer := facade.Lang().GetLocalizer("en")
func Lang() *LangFacade {
	return &LangFacade{}
}

type LangFacade struct{}

// T 翻译
func (l *LangFacade) T(ctx context.Context, messageID string, data map[string]interface{}) string {
	return lang.T(ctx, messageID, data)
}

// GetLocalizer 获取指定语言的Localizer
func (l *LangFacade) GetLocalizer(langCode string) *i18n.Localizer {
	return lang.GetLocalizer(langCode)
}

// GetBundle 获取翻译包
func (l *LangFacade) GetBundle() *i18n.Bundle {
	return lang.GetBundle()
}

// IsLoaded 检查翻译是否已加载
func (l *LangFacade) IsLoaded() bool {
	return lang.IsLoaded()
}
