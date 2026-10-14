package middleware

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/common/ctxkey"
	"gin/pkg/serviceprovider/debugger"
	"strings"
	"time"
	"uuid"

	"github.com/gin-gonic/gin"
)

type Logger struct {
	base.BaseMiddleware
}

// Handle 日志中间件
func (s Logger) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceId := uuid.New().String()
		lang := s.getLang(c)

		// 构建完整的context
		ctx := c.Request.Context()

		// 注入追踪信息
		ctx = ctxkey.WithValue(ctx, ctxkey.TraceIDKey, traceId)
		ctx = ctxkey.WithValue(ctx, ctxkey.IpKey, c.ClientIP())
		ctx = ctxkey.WithValue(ctx, ctxkey.PathKey, c.Request.URL.Path)
		ctx = ctxkey.WithValue(ctx, ctxkey.MethodKey, c.Request.Method)
		ctx = ctxkey.WithValue(ctx, ctxkey.ParamsKey, facade.Request().Params(c))
		ctx = ctxkey.WithValue(ctx, ctxkey.LangKey, lang)
		ctx = ctxkey.WithValue(ctx, ctxkey.StartTimeKey, start)

		c.Header("Trace-Id", traceId)

		// 注入日志
		ctx = ctxkey.WithValue(ctx, ctxkey.DebuggerKey, facade.Log().WithDebugger(ctx))

		// 更新请求的context
		c.Request = c.Request.WithContext(ctx)

		defer func() {
			// 从当前请求获取最新的ctx
			currentCtx := c.Request.Context()
			cost := float64(time.Since(start).Milliseconds())

			// 更新耗时
			currentCtx = ctxkey.WithValue(currentCtx, ctxkey.MsKey, cost)
			c.Request = c.Request.WithContext(currentCtx)

			if facade.Config().Log.Access {
				facade.Log().WithDebugger(currentCtx).Info("Access Log")
			}
			debugger.Store.Delete(traceId)
		}()

		// 执行后续处理
		c.Next()
	}
}

// getLang 获取语言
func (s Logger) getLang(c *gin.Context) string {
	// 配置支持的语言如["zh", "en"]
	supported := strings.Split(facade.Config().I18n.Lang, ",")

	if q := strings.ToLower(strings.TrimSpace(c.Query("lang"))); q != "" {
		if lang := s.matchLang(q, supported); lang != "" {
			return lang
		}
	}

	if h := strings.ToLower(c.GetHeader("Accept-Language")); h != "" {
		if lang := s.matchLang(h, supported); lang != "" {
			return lang
		}
	}

	return "zh"
}

// matchLang 匹配语言
func (s Logger) matchLang(input string, supported []string) string {
	for _, lang := range supported {
		lang = strings.ToLower(strings.TrimSpace(lang))
		// 支持en/en-US/en_US/zh/zh-CN/zh_CN
		if strings.HasPrefix(input, lang) {
			return lang
		}
	}
	return ""
}
