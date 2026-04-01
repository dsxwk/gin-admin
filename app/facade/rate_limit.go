package facade

import (
	"context"
	"gin/common/errcode"
	"gin/common/response"
	"gin/pkg/ratelimit"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

var (
	rlErr         = errcode.RateLimitError()
	userStore     *ratelimit.Store
	ipStore       *ratelimit.Store
	rateLimitOnce sync.Once
)

// initRateLimit 初始化限流
func initRateLimit() {
	rateLimitOnce.Do(func() {
		userStore = ratelimit.NewStore(5*time.Minute, 100, 200)
		ipStore = ratelimit.NewStore(5*time.Minute, 100, 200)
	})
}

// shutdownRateLimit 关闭限流
func shutdownRateLimit() {
	if userStore != nil {
		userStore.Close()
	}
	if ipStore != nil {
		ipStore.Close()
	}
}

// globalRateLimit 全局限流中间件
func globalRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if userStore == nil {
			c.Next()
			return
		}
		if !userStore.AllowGlobal() {
			response.Error(c, &rlErr)
			return
		}
		c.Next()
	}
}

// ipRateLimit IP限流中间件
func ipRateLimit(r rate.Limit, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if ipStore == nil {
			c.Next()
			return
		}
		key := c.ClientIP() + ":" + c.FullPath()
		if !ipStore.AllowKey(key, r, burst) {
			response.Error(c, &rlErr)
			return
		}
		c.Next()
	}
}

// userRateLimit 用户限流中间件
func userRateLimit(r rate.Limit, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user.id")
		if userID == "" {
			c.Next()
			return
		}
		key := userID + ":" + c.FullPath()
		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Millisecond)
		defer cancel()
		if err := userStore.WaitKey(ctx, key, r, burst); err != nil {
			response.Error(c, &rlErr)
			return
		}
		c.Next()
	}
}

// RateLimiter 限流门面
var RateLimiter = &rateLimiterFacade{}

type rateLimiterFacade struct{}

// Global 全局限流中间件
func (rl *rateLimiterFacade) Global() gin.HandlerFunc {
	return globalRateLimit()
}

// IP IP限流中间件
// r: 每秒产生多少token
// burst: 桶容量
func (rl *rateLimiterFacade) IP(r rate.Limit, burst int) gin.HandlerFunc {
	return ipRateLimit(r, burst)
}

// User 用户限流中间件
// r: 每秒产生多少token
// burst: 桶容量
func (rl *rateLimiterFacade) User(r rate.Limit, burst int) gin.HandlerFunc {
	return userRateLimit(r, burst)
}

// Init 初始化限流器
func (rl *rateLimiterFacade) Init() {
	initRateLimit()
}

// Shutdown 关闭限流器
func (rl *rateLimiterFacade) Shutdown() {
	shutdownRateLimit()
}
