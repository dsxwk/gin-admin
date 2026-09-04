package base

import (
	"context"
	"gin/app/facade"
)

type BaseRequest struct {
	Context
}

func (s *BaseRequest) WithContext(ctx context.Context) *BaseRequest {
	s.SetContext(ctx)

	return s
}

func (s *BaseRequest) Trans(ctx context.Context, messageID string, data map[string]any) string {
	return facade.Lang().Trans(ctx, messageID, data)
}
