package middleware

import (
	"gin/common/base"
	"gin/common/errcode"
	"gin/common/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"time"
)

var (
	rlErr = errcode.RateLimitError()
	// 限流存储5分钟不访问删除
	userStore = newLimiterStore(5 * time.Minute)
	ipStore   = newLimiterStore(5 * time.Minute)
)

type RateLimit struct {
	base.BaseMiddleware
	limiter *rate.Limiter
}

// NewRateLimit 创建限流中间件
// r 每秒产生多少token
// burst 桶容量
func NewRateLimit(r rate.Limit, burst int) *RateLimit {
	return &RateLimit{
		limiter: rate.NewLimiter(r, burst),
	}
}

// Handle 限流中间件
func (s RateLimit) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !s.limiter.Allow() {
			response.Error(c, &rlErr)
			return
		}
		c.Next()
	}
}

// UserRateLimit 用户限流
// r 每秒产生多少token
// burst 桶容量
func UserRateLimit(r rate.Limit, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user.id")
		if userID == "" {
			c.Next()
			return
		}
		limiter := userStore.get(userID, r, burst)

		if !limiter.Allow() {
			response.Error(c, &rlErr)
			return
		}

		c.Next()
	}
}

// IpRateLimit ip限流
// r 每秒产生多少token
// burst 桶容量
func IpRateLimit(r rate.Limit, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := ipStore.get(ip, r, burst)

		if !limiter.Allow() {
			response.Error(c, &rlErr)
			return
		}

		c.Next()
	}
}
