package middleware

import (
	"gin/app/facade"
	"gin/common/base"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimit struct {
	base.BaseMiddleware
}

// Handle 全局限流
func (r RateLimit) Handle() gin.HandlerFunc {
	return facade.RateLimiter.Global()
}

// IpRateLimit IP限流
func (r RateLimit) IpRateLimit(rps rate.Limit, burst int) gin.HandlerFunc {
	return facade.RateLimiter.IP(rps, burst)
}

// UserRateLimit 用户限流
func (r RateLimit) UserRateLimit(rps rate.Limit, burst int) gin.HandlerFunc {
	return facade.RateLimiter.User(rps, burst)
}
