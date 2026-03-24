package middleware

import (
	"context"
	"gin/common/base"
	"gin/common/errcode"
	"gin/common/response"
	"gin/pkg/ratelimit"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

var (
	rlErr     = errcode.RateLimitError()
	userStore *ratelimit.Store
	ipStore   *ratelimit.Store
	once      sync.Once
)

type RateLimit struct {
	base.BaseMiddleware
	limiter *rate.Limiter
}

// InitRateLimit 初始化限流(必须在main启动时调用)
func InitRateLimit() {
	// 限流存储5分钟不访问删除
	once.Do(func() {
		// 全局限流: 100QPS,突发200
		userStore = ratelimit.NewStore(5*time.Minute, 100, 200)
		ipStore = ratelimit.NewStore(5*time.Minute, 100, 200)
	})
}

// ShutdownRateLimit 关闭限流(必须调用)
func ShutdownRateLimit() {
	if userStore != nil {
		userStore.Close()
	}
	if ipStore != nil {
		ipStore.Close()
	}
}

// Handle 全局限流
func (s RateLimit) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防御: 未初始化
		if userStore == nil {
			c.Next()
			return
		}

		// 全局令牌桶
		if !userStore.AllowGlobal() {
			response.Error(c, &rlErr)
			return
		}

		c.Next()
	}
}

// UserRateLimit 用户限流
// r 每秒产生多少token
// burst 桶容量
func (s RateLimit) UserRateLimit(r rate.Limit, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user.id")
		if userID == "" {
			c.Next()
			return
		}
		// 用户+路径(接口级限流)
		key := userID + ":" + c.FullPath()

		// 平滑限流(最多等待100ms)
		ctx, cancel := context.WithTimeout(c.Request.Context(), 100*time.Millisecond)
		defer cancel()

		if err := userStore.WaitKey(ctx, key, r, burst); err != nil {
			response.Error(c, &rlErr)
			return
		}

		c.Next()
	}
}

// IpRateLimit ip限流
// r 每秒产生多少token
// burst 桶容量
func (s RateLimit) IpRateLimit(r rate.Limit, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if ipStore == nil {
			c.Next()
			return
		}

		ip := c.ClientIP()
		// ip+路径
		key := ip + ":" + c.FullPath()

		// 快速失败
		if !ipStore.AllowKey(key, r, burst) {
			response.Error(c, &rlErr)
			return
		}

		c.Next()
	}
}
