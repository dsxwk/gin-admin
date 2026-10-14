package base

import (
	"context"
)

// Context 统一管理context
type Context struct {
	context context.Context `swaggerignore:"true"`
}

// WithContext 设置上下文
func (s *Context) WithContext(ctx context.Context) {
	s.SetContext(ctx)
}

// SetContext 设置上下文
func (s *Context) SetContext(ctx context.Context) {
	s.context = ctx
}

// Context 获取ctx
func (s *Context) Context() context.Context {
	return s.context
}
