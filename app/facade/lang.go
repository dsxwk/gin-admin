package facade

import (
	"context"
	"gin/pkg/container"
	"gin/pkg/serviceprovider/lang"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Lang 翻译门面
// 使用示例:
//
//	msg := facade.Lang().Trans(ctx, "welcome", map[string]any{"name": "John"})
//	localizer := facade.Lang().Localizer("en")
func Lang() *LangFacade {
	return &LangFacade{
		service: container.Default().Lang(),
	}
}

// LangFacade 翻译门面
type LangFacade struct {
	service *lang.Service
}

// Trans 翻译
func (l *LangFacade) Trans(ctx context.Context, messageID string, data map[string]any) string {
	if l.service == nil {
		return messageID
	}
	return l.service.Trans(ctx, messageID, data)
}

// Localizer 获取指定语言的Localizer
func (l *LangFacade) Localizer(langCode string) *i18n.Localizer {
	if l.service == nil {
		return nil
	}
	return l.service.Localizer(langCode)
}

// Bundle 获取翻译包
func (l *LangFacade) Bundle() *i18n.Bundle {
	if l.service == nil {
		return nil
	}
	return l.service.Bundle()
}

// IsLoaded 检查翻译是否已加载
func (l *LangFacade) IsLoaded() bool {
	if l.service == nil {
		return false
	}
	return l.service.IsLoaded()
}
