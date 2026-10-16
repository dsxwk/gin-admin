package base

import (
	"context"
	"gin/app/facade"
	"gin/pkg/errcode"
)

// Middleware 基础中间件
type Middleware struct {
	Response errcode.Response
}

// Trans 翻译消息
func (s *Middleware) Trans(ctx context.Context, messageID string, data map[string]any) string {
	return facade.Lang().Trans(ctx, messageID, data)
}
