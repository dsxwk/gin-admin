package ratelimit

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

const (
	// ipKeyPrefix ip限流key前缀
	ipKeyPrefix = "ip:"
	// userKeyPrefix 用户限流key前缀
	userKeyPrefix = "user:"
)

// Manager 限流管理器
type Manager struct {
	global *rate.Limiter // 全局令牌桶
	store  *keyedStore   // key级令牌桶存储
}

// NewManager 创建限流管理器
func NewManager(ttl time.Duration, globalRate rate.Limit, globalBurst int) *Manager {
	return &Manager{
		global: rate.NewLimiter(globalRate, globalBurst),
		store:  newKeyedStore(ttl),
	}
}

// AllowGlobal 全局限流
func (m *Manager) AllowGlobal() bool {
	if m == nil || m.global == nil {
		return true
	}
	return m.global.Allow()
}

// AllowIP IP限流
func (m *Manager) AllowIP(ip, path string, r rate.Limit, burst int) bool {
	if m == nil || m.store == nil {
		return true
	}
	return m.store.allow(ipKeyPrefix+ip+":"+path, r, burst)
}

// WaitUser 用户限流
func (m *Manager) WaitUser(ctx context.Context, userID, path string, r rate.Limit, burst int) error {
	if m == nil || m.store == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return m.store.wait(ctx, userKeyPrefix+userID+":"+path, r, burst)
}

// Close 关闭限流管理器
func (m *Manager) Close() {
	if m == nil || m.store == nil {
		return
	}
	m.store.close()
}
