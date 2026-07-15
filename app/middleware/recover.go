package middleware

import (
	"context"
	"fmt"
	"gin/app/facade"
	"gin/common/base"
	"gin/common/ctxkey"
	"gin/common/errcode"
	"github.com/gin-gonic/gin"
	"runtime"
)

type Recover struct {
	base.BaseMiddleware
}

type ErrData struct {
	TraceId string      `json:"traceId"`
	Error   interface{} `json:"error"`
	IP      string      `json:"ip"`
	Lang    string      `json:"lang"`
	Path    string      `json:"path"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Stack   []string    `json:"stack"`
}

// Handle recover中间件
func (s Recover) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx := c.Request.Context()
				stack := getStackTrace(3)

				if facade.Config().App.Env == "production" {
					s.Response.Error(c, errcode.SystemError().WithMsg(fmt.Sprintf("%v", err)))
				} else {
					s.Response.Error(
						c,
						errcode.SystemError().
							WithMsg(fmt.Sprintf("%v", err)).
							WithData(
								&ErrData{
									TraceId: getString(ctx, ctxkey.TraceIdKey),
									Error:   err,
									IP:      getString(ctx, ctxkey.IpKey),
									Lang:    getString(ctx, ctxkey.LangKey),
									Path:    getString(ctx, ctxkey.PathKey),
									Method:  getString(ctx, ctxkey.MethodKey),
									Params:  ctx.Value(ctxkey.ParamsKey),
									Stack:   stack,
								},
							),
					)
				}
			}
		}()
		c.Next()
	}
}

func getString(c context.Context, key string) string {
	if c == nil {
		return "unknown"
	}
	if v, ok := c.Value(key).(string); ok {
		return v
	}
	return "unknown"
}

func getStackTrace(skip int) []string {
	const maxDepth = 32
	pc := make([]uintptr, maxDepth)
	n := runtime.Callers(skip, pc)
	pc = pc[:n]

	var trace []string
	for _, p := range pc {
		fn := runtime.FuncForPC(p)
		if fn == nil {
			trace = append(trace, "unknown")
			continue
		}
		file, line := fn.FileLine(p)
		trace = append(trace, fmt.Sprintf("%s:%d %s", file, line, fn.Name()))
	}

	return trace
}
