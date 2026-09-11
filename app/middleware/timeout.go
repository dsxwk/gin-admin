package middleware

import (
	"context"
	"errors"
	"gin/app/errcode"
	"gin/common/base"
	"time"

	"github.com/gin-gonic/gin"
)

const SkipTimeoutKey = "_skip_timeout"

type Timeout struct {
	base.BaseMiddleware
}

// timeoutWriter 超时感知的ResponseWriter
type timeoutWriter struct {
	gin.ResponseWriter
	ctx      context.Context
	timedOut bool
}

// expired 判断请求是否已超时
func (w *timeoutWriter) expired() bool {
	return w.timedOut || errors.Is(w.ctx.Err(), context.DeadlineExceeded)
}

func (w *timeoutWriter) Write(data []byte) (int, error) {
	if w.expired() {
		w.timedOut = true
		return 0, errors.New("request timed out")
	}
	return w.ResponseWriter.Write(data)
}

func (w *timeoutWriter) WriteHeader(statusCode int) {
	if w.expired() {
		w.timedOut = true
		return
	}
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *timeoutWriter) WriteString(s string) (int, error) {
	if w.expired() {
		w.timedOut = true
		return 0, errors.New("request timed out")
	}
	return w.ResponseWriter.WriteString(s)
}

// Handle 超时中间件
func (s Timeout) Handle(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果标记了跳过,直接放行
		if _, ok := c.Get(SkipTimeoutKey); ok {
			c.Next()
			return
		}

		if timeout <= 0 {
			c.Next()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		// 包装ResponseWriter,超时后阻止后续写入
		originalWriter := c.Writer
		writer := &timeoutWriter{ResponseWriter: originalWriter, ctx: ctx}
		c.Writer = writer

		c.Next()

		// 检查是否超时
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			writer.timedOut = true
			c.Writer = originalWriter
			// 如果还没写入响应,写超时响应
			if !originalWriter.Written() {
				c.Abort()
				s.Response.Error(c, errcode.TimeoutError())
			}
		}
	}
}
