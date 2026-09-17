package middleware

import (
	"context"
	"gin/app/errcode"
	"gin/app/facade"
	"gin/common/base"
	"gin/common/ctxkey"
	"gin/pkg"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimit struct {
	base.BaseMiddleware
}

// Handle 全局限流
func (r RateLimit) Handle() gin.HandlerFunc {
	limiter := facade.RateLimiter()
	return func(c *gin.Context) {
		if limiter == nil || limiter.AllowGlobal() {
			c.Next()
			return
		}
		facade.Response().Error(c, errcode.RateLimitError())
	}
}

// IpRateLimit IP限流
func (r RateLimit) IpRateLimit(rps rate.Limit, burst int) gin.HandlerFunc {
	limiter := facade.RateLimiter()
	return func(c *gin.Context) {
		if limiter == nil || limiter.AllowIP(c.ClientIP(), c.FullPath(), rps, burst) {
			c.Next()
			return
		}
		facade.Response().Error(c, errcode.RateLimitError())
	}
}

// UserRateLimit 用户限流
func (r RateLimit) UserRateLimit(rps rate.Limit, burst int) gin.HandlerFunc {
	limiter := facade.RateLimiter()
	return func(c *gin.Context) {
		if limiter == nil {
			c.Next()
			return
		}

		userID := c.GetInt64(ctxkey.UserIdKey)
		if userID <= 0 {
			c.Next()
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Millisecond)
		defer cancel()
		if err := limiter.WaitUser(ctx, pkg.IntToString[int64](userID), c.FullPath(), rps, burst); err != nil {
			facade.Response().Error(c, errcode.RateLimitError())
			return
		}
		c.Next()
	}
}
