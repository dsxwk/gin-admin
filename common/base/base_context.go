package base

import (
	"context"
	"gin/common/ctxkey"
)

// Context 统一管理context
type Context struct {
	Ctx context.Context
}

// Set 设置ctx
func (s *Context) Set(ctx context.Context) {
	s.Ctx = ctx
}

// Get 获取ctx
func (s *Context) Get() context.Context {
	return s.Ctx
}

// TraceId 获取traceId
func (s *Context) TraceId() string {
	return getString(s.Ctx, ctxkey.TraceIdKey)
}

// GetLang 获取语言
func (s *Context) GetLang() string {
	return getString(s.Ctx, ctxkey.LangKey)
}

// GetIp 获取ip
func (s *Context) GetIp() string {
	return getString(s.Ctx, ctxkey.IpKey)
}

// GetPath 获取请求路径
func (s *Context) GetPath() string {
	return getString(s.Ctx, ctxkey.PathKey)
}

// GetMethod 获取请求方法
func (s *Context) GetMethod() string {
	return getString(s.Ctx, ctxkey.MethodKey)
}

// GetParams 获取请求参数
func (s *Context) GetParams() string {
	return getString(s.Ctx, ctxkey.ParamsKey)
}

// GetMs 获取耗时
func (s *Context) GetMs() string {
	return getString(s.Ctx, ctxkey.MsKey)
}

// GetStartTime 获取请求开始时间
func (s *Context) GetStartTime() string {
	return getString(s.Ctx, ctxkey.StartTimeKey)
}

// 防止panic
func getString(c context.Context, key string) string {
	if v, ok := c.Value(key).(string); ok {
		return v
	}
	return "unknown"
}
