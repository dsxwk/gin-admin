package ctxkey

import "context"

const (
	UserIdKey    string = "userId"
	TraceIDKey   string = "traceId"
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

// Value 从context获取值
func Value(ctx context.Context, key string) any {
	return ctx.Value(key)
}

// TraceID 获取tracId
func TraceID(ctx context.Context) string {
	if ctx == nil {
		return "unknown"
	}
	if id := ctx.Value(TraceIDKey); id != nil {
		if s, ok := id.(string); ok && s != "" {
			return s
		}
	}
	return "unknown"
}
