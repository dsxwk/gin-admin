package middleware

import (
	"bytes"
	"gin/app/facade"
	"gin/common/base"
	"gin/common/ctxkey"
	"gin/pkg/debugger"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"io"
	"net/http"
	"strings"
	"time"
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
		ctx = ctxkey.WithValue(ctx, ctxkey.TraceIdKey, traceId)
		ctx = ctxkey.WithValue(ctx, ctxkey.IpKey, c.ClientIP())
		ctx = ctxkey.WithValue(ctx, ctxkey.PathKey, c.Request.URL.Path)
		ctx = ctxkey.WithValue(ctx, ctxkey.MethodKey, c.Request.Method)
		ctx = ctxkey.WithValue(ctx, ctxkey.ParamsKey, s.getParams(c))
		ctx = ctxkey.WithValue(ctx, ctxkey.LangKey, lang)
		ctx = ctxkey.WithValue(ctx, ctxkey.StartTimeKey, start)

		c.Header("Trace-Id", traceId)

		// 注入日志
		log := facade.Log.Logger()
		ctx = ctxkey.WithValue(ctx, ctxkey.DebuggerKey, log.WithDebugger(ctx))

		// 更新请求的context
		c.Request = c.Request.WithContext(ctx)

		defer func() {
			// 从当前请求获取最新的ctx
			currentCtx := c.Request.Context()
			cost := float64(time.Since(start).Milliseconds())

			// 更新耗时
			currentCtx = ctxkey.WithValue(currentCtx, ctxkey.MsKey, cost)
			c.Request = c.Request.WithContext(currentCtx)

			if conf.Log.Access {
				log.WithDebugger(currentCtx).Info("Access Log")
			}
			debugger.Store.Delete(traceId)
		}()

		// 执行后续处理
		c.Next()
	}
}

// 获取参数
func (s Logger) getParams(c *gin.Context) any {
	// GET/DELETE/query 参数
	if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodDelete {
		return c.Request.URL.Query()
	}

	// 其他方法尝试读取body
	if c.Request.Body == nil {
		return map[string]any{}
	}

	body, _err := io.ReadAll(c.Request.Body)
	if _err != nil || len(body) == 0 {
		return map[string]any{}
	}

	// 读完要塞回去,避免后续handler读不到
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	// 尝试解析json
	var m map[string]any
	if _e := json.Unmarshal(body, &m); _e == nil {
		return m
	}

	// 非json原样记录
	return string(body)
}

// getLang 获取语言
func (s *Logger) getLang(c *gin.Context) string {
	// 配置支持的语言如["zh", "en"]
	supported := strings.Split(conf.I18n.Lang, ",")

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
func (s *Logger) matchLang(input string, supported []string) string {
	for _, lang := range supported {
		lang = strings.ToLower(strings.TrimSpace(lang))
		// 支持en/en-US/en_US/zh/zh-CN/zh_CN
		if strings.HasPrefix(input, lang) {
			return lang
		}
	}
	return ""
}
