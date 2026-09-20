package base

import (
	"gin/app/facade"
	"gin/pkg/errcode"
)

type BaseMiddleware struct {
	Context
	Response errcode.Response
}

func (s *BaseMiddleware) Trans(messageID string, data map[string]any) string {
	return facade.Lang().Trans(s.Context.Context(), messageID, data)
}
