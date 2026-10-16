package ctxkey

import (
	"context"
	"time"
)

// Key 标准Context键类型
type Key uint8

const (
	// UserIdKey 用户ID
	UserIdKey Key = iota
	// TraceIDKey 链路ID
	TraceIDKey
	// IpKey IP地址
	IpKey
	// PathKey 请求路径
	PathKey
	// MethodKey 请求方法
	MethodKey
	// ParamsKey 请求参数
	ParamsKey
	// MsKey 请求耗时
	MsKey
	// LangKey 语言
	LangKey
	// StartTimeKey 开始时间
	StartTimeKey
	// DebuggerKey 调试器
	DebuggerKey
)

const (
	// GinUserIDKey Gin上下文用户ID
	GinUserIDKey string = "userId"
)

// WithValue 将值注入到标准Context
func WithValue(ctx context.Context, key Key, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

// Value 从标准Context获取值
func Value(ctx context.Context, key Key) any {
	if ctx == nil {
		return nil
	}
	return ctx.Value(key)
}

// TraceID 获取链路ID
func TraceID(ctx context.Context) string {
	id, _ := Value(ctx, TraceIDKey).(string)
	if id == "" {
		return "unknown"
	}
	return id
}

// UserID 获取用户ID
func UserID(ctx context.Context) int64 {
	switch id := Value(ctx, UserIdKey).(type) {
	case int64:
		return id
	case int:
		return int64(id)
	case float64:
		return int64(id)
	default:
		return 0
	}
}

// Lang 获取语言
func Lang(ctx context.Context) string {
	lang, _ := Value(ctx, LangKey).(string)
	return lang
}

// StartTime 获取开始时间
func StartTime(ctx context.Context) time.Time {
	startTime, _ := Value(ctx, StartTimeKey).(time.Time)
	return startTime
}

// Ms 获取请求耗时
func Ms(ctx context.Context) float64 {
	ms, _ := Value(ctx, MsKey).(float64)
	return ms
}
