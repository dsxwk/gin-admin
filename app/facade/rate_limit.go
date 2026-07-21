package facade

import (
	"context"
	"gin/common/ctxkey"
	"gin/common/errcode"
	"gin/common/response"
	"gin/pkg"
	"gin/pkg/serviceprovider/ratelimit"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"time"
)

var (
	userStore *ratelimit.Store
	ipStore   *ratelimit.Store
)

// initRateLimit 初始化限流
func initRateLimit() {
	if userStore != nil {
		return
	}
	userStore = ratelimit.NewStore(5*time.Minute, 100, 200)
	ipStore = ratelimit.NewStore(5*time.Minute, 100, 200)
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
			response.Response{}.Error(c, errcode.RateLimitError())
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
			response.Response{}.Error(c, errcode.RateLimitError())
			return
		}
		c.Next()
	}
}

// userRateLimit 用户限流中间件
func userRateLimit(r rate.Limit, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID string
		if id := c.GetInt64(ctxkey.UserIdKey); id > 0 {
			userID = pkg.IntToString[int64](id)
		}
		if userID == "" {
			c.Next()
			return
		}
		key := userID + ":" + c.FullPath()
		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Millisecond)
		defer cancel()
		if err := userStore.WaitKey(ctx, key, r, burst); err != nil {
			response.Response{}.Error(c, errcode.RateLimitError())
			return
		}
		c.Next()
	}
}

// RateLimiter 限流门面
// 使用示例:
//
//	router.Use(facade.RateLimiter().Global())
//	router.Use(facade.RateLimiter().IP(10, 20))
func RateLimiter() *RateLimiterFacade {
	return &RateLimiterFacade{}
}

type RateLimiterFacade struct{}

// Global 全局限流中间件
func (rl *RateLimiterFacade) Global() gin.HandlerFunc {
	return globalRateLimit()
}

// IP IP限流中间件
// r: 每秒产生多少token
// burst: 桶容量
func (rl *RateLimiterFacade) IP(r rate.Limit, burst int) gin.HandlerFunc {
	return ipRateLimit(r, burst)
}

// User 用户限流中间件
// r: 每秒产生多少token
// burst: 桶容量
func (rl *RateLimiterFacade) User(r rate.Limit, burst int) gin.HandlerFunc {
	return userRateLimit(r, burst)
}

// Init 初始化限流器
func (rl *RateLimiterFacade) Init() {
	initRateLimit()
}

// Shutdown 关闭限流器
func (rl *RateLimiterFacade) Shutdown() {
	shutdownRateLimit()
}
