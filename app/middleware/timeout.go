package middleware

import (
	"context"
	"gin/common/base"
	"gin/common/errcode"
	"gin/common/response"
	"github.com/gin-gonic/gin"
	"time"
)

type Timeout struct {
	base.BaseMiddleware
}

// Handle 超时中间件
func (s Timeout) Handle(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})

		go func() {
			c.Next()
			close(done)
		}()

		select {

		case <-done:
			return

		case <-ctx.Done():
			errCode := errcode.TimeoutError()
			response.Error(c, &errCode)
			return
		}
	}
}
