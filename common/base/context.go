package base

import (
	"context"
)

// Context 统一管理上下文
type Context struct {
	context context.Context `swaggerignore:"true"`
}

// SetContext 设置上下文
func (s *Context) SetContext(ctx context.Context) {
	s.context = ctx
}

// GetContext 获取上下文
func (s *Context) GetContext() context.Context {
	return s.context
}
