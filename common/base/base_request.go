package base

import (
	"gin/app/facade"
)

type BaseRequest struct {
	Context
}

func (s *BaseRequest) Trans(messageID string, data map[string]any) string {
	return facade.Lang().Trans(s.Context.Context(), messageID, data)
}
