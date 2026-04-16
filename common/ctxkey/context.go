package ctxkey

import "context"

type contextKey string

const (
	UserIdKey    string     = "userId"
	TraceIdKey   contextKey = "traceId"
	IpKey        contextKey = "ip"
	PathKey      contextKey = "path"
	MethodKey    contextKey = "method"
	ParamsKey    contextKey = "params"
	MsKey        contextKey = "ms"
	LangKey      contextKey = "lang"
	StartTimeKey contextKey = "startTime"
	DebuggerKey  contextKey = "debugger"
)

// WithValue 将值注入到context
func WithValue(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

// GetValue 从context获取值
func GetValue(ctx context.Context, key string) interface{} {
	return ctx.Value(key)
}
