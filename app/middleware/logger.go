package middleware

import (
	"bytes"
	"context"
	"gin/common/base"
	"gin/common/ctxkey"
	"gin/common/trace"
	"gin/config"
	"gin/pkg/container"
	"gin/pkg/logger"
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
		lang := s.GetLang(c)
		ctn := container.GetContainer()

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, ctxkey.TraceIdKey, traceId)
		ctx = context.WithValue(ctx, ctxkey.IpKey, c.ClientIP())
		ctx = context.WithValue(ctx, ctxkey.PathKey, c.Request.URL.Path)
		ctx = context.WithValue(ctx, ctxkey.MethodKey, c.Request.Method)
		ctx = context.WithValue(ctx, ctxkey.ParamsKey, s.getParams(c))
		ctx = context.WithValue(ctx, ctxkey.LangKey, lang)
		ctx = context.WithValue(ctx, ctxkey.StartTimeKey, start)

		ctnCtx := ctn.WithContext(ctx)
		ctx = container.Set(ctx, ctnCtx)

		c.Request = c.Request.WithContext(ctx)
		c.Header("Trace-Id", traceId)

		defer func() {
			// 清理trace
			trace.Store.Delete(traceId)
			cost := float64(time.Since(start).Milliseconds())
			ctx = context.WithValue(ctx, ctxkey.MsKey, cost)
			c.Request = c.Request.WithContext(ctx)
		}()

		c.Next()

		if config.NewConfig().Log.Access {
			logger.NewLogger().WithDebugger(ctx).Info("Access Log")
		}
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

// GetLang 获取语言
func (s *Logger) GetLang(c *gin.Context) string {
	// 配置支持的语言如["zh", "en"]
	supported := strings.Split(config.NewConfig().I18n.Lang, ",")

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
