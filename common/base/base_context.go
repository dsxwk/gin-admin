package base

import (
	"context"
	"gin/common/ctxkey"
)

// Context 统一管理context
type Context struct {
	context context.Context `swaggerignore:"true"`
}

// WithContext 设置上下文
func (s *Context) WithContext(ctx context.Context) {
	s.context = ctx
}

// Context 获取ctx
func (s *Context) Context() context.Context {
	return s.context
}

// TraceID 获取traceId
func (s *Context) TraceID() string {
	return getString(s.Context(), ctxkey.TraceIDKey)
}

// GetLang 获取语言
func (s *Context) GetLang() string {
	return getString(s.Context(), ctxkey.LangKey)
}

// GetIp 获取ip
func (s *Context) GetIp() string {
	return getString(s.Context(), ctxkey.IpKey)
}

// GetPath 获取请求路径
func (s *Context) GetPath() string {
	return getString(s.Context(), ctxkey.PathKey)
}

// GetMethod 获取请求方法
func (s *Context) GetMethod() string {
	return getString(s.Context(), ctxkey.MethodKey)
}

// GetParams 获取请求参数
func (s *Context) GetParams() string {
	return getString(s.Context(), ctxkey.ParamsKey)
}

// GetMs 获取耗时
func (s *Context) GetMs() string {
	return getString(s.Context(), ctxkey.MsKey)
}

// GetStartTime 获取请求开始时间
func (s *Context) GetStartTime() string {
	return getString(s.Context(), ctxkey.StartTimeKey)
}

// GetUserId 获取当前用户ID
func (s *Context) GetUserId() int64 {
	if s.Context() == nil {
		return 0
	}
	if v, ok := s.Context().Value(ctxkey.UserIdKey).(int64); ok {
		return v
	}
	return 0
}

// 防止panic
func getString(c context.Context, key string) string {
	if c == nil {
		return "unknown"
	}
	if v, ok := c.Value(key).(string); ok {
		return v
	}
	return "unknown"
}
