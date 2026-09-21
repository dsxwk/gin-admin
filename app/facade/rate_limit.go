package facade

import (
	"gin/pkg/container"
	"gin/pkg/serviceprovider/ratelimit"
)

// RateLimiter 获取限流管理器
func RateLimiter() *ratelimit.Manager {
	return container.Default().RateLimit()
}
