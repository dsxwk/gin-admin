package base

import "context"

type BaseRequest struct {
	Context
}

func (s *BaseRequest) WithContext(ctx context.Context) *BaseRequest {
	s.Set(ctx)

	return s
}
