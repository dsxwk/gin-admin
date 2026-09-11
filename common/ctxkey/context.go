package ctxkey

import "context"

const (
	UserIdKey    string = "userId"
	TraceIdKey   string = "traceId"
	IpKey        string = "ip"
	PathKey      string = "path"
	MethodKey    string = "method"
	ParamsKey    string = "params"
	MsKey        string = "ms"
	LangKey      string = "lang"
	StartTimeKey string = "startTime"
	DebuggerKey  string = "debugger"
)

// WithValue 将值注入到context
func WithValue(ctx context.Context, key string, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

// GetValue 从context获取值
func GetValue(ctx context.Context, key string) any {
	return ctx.Value(key)
}

// GetTraceId 获取tracId
func GetTraceId(ctx context.Context) string {
	if ctx == nil {
		return "unknown"
	}
	if id := ctx.Value(TraceIdKey); id != nil {
		if s, ok := id.(string); ok && s != "" {
			return s
		}
	}
	return "unknown"
}
