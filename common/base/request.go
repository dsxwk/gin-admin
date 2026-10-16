package base

import (
	"gin/app/facade"
)

// Request 基础请求
type Request struct {
	Context
}

// Trans 翻译消息
func (s *Request) Trans(messageID string, data map[string]any) string {
	return facade.Lang().Trans(s.GetContext(), messageID, data)
}
