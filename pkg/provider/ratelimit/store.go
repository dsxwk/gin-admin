package ratelimit

import (
	"context"
	"golang.org/x/time/rate"
	"sync"
	"sync/atomic"
	"time"
)

// 每个key对应一个令牌桶
type limiterItem struct {
	limiter  *rate.Limiter // 令牌桶
	lastSeen int64         // 最后访问时间(用于TTL清理,必须原子操作)
	r        rate.Limit    // 当前速率(用于动态更新)
	burst    int           // 桶容量
}

// Store 限流存储容器
type Store struct {
	m      sync.Map       // key->limiterItem(高并发读写)
	ttl    time.Duration  // 多久未访问自动删除
	stop   chan struct{}  // 用于停止清理协程
	wg     sync.WaitGroup // 等待goroutine退出
	global *rate.Limiter  // 全局令牌桶
}

// NewStore 创建限流存储(带自动清理)
func NewStore(ttl time.Duration, r rate.Limit, burst int) *Store {
	s := &Store{
		ttl:    ttl,
		stop:   make(chan struct{}),
		global: rate.NewLimiter(r, burst),
	}

	// 标记一个goroutine
	s.wg.Add(1)

	// 启动后台清理协程
	go s.clean()

	return s
}

// Close 关闭清理协程(必须在服务关闭时调用)
func (s *Store) Close() {
	select {
	case <-s.stop:
		// 已关闭,防止重复close panic
		return
	default:
		close(s.stop)
	}

	// 等待clean goroutine退出
	s.wg.Wait()
}

// AllowGlobal 全局限流(快速失败)
// return false=超过系统QPS限制
func (s *Store) AllowGlobal() bool {
	return s.global.Allow()
}

// WaitGlobal 全局限流(平滑限流)
// 会等待一段时间获取token
func (s *Store) WaitGlobal(ctx context.Context) error {
	return s.global.Wait(ctx)
}

// Get 获取某个key对应的令牌桶
// 支持: 自动创建,自动复用,参数变化自动更新
func (s *Store) Get(key string, r rate.Limit, burst int) *rate.Limiter {
	now := time.Now().Unix()

	// 快路径(已有)
	if v, ok := s.m.Load(key); ok {
		item := v.(*limiterItem)

		// 如果限流参数变了→重新创建limiter
		if item.r != r || item.burst != burst {
			newItem := &limiterItem{
				limiter:  rate.NewLimiter(r, burst),
				lastSeen: now,
				r:        r,
				burst:    burst,
			}
			s.m.Store(key, newItem)
			return newItem.limiter
		}

		// 更新时间(必须原子)
		atomic.StoreInt64(&item.lastSeen, now)
		return item.limiter
	}

	// 慢路径(首次创建)
	item := &limiterItem{
		limiter:  rate.NewLimiter(r, burst),
		lastSeen: now,
		r:        r,
		burst:    burst,
	}

	// 防止并发重复创建
	actual, _ := s.m.LoadOrStore(key, item)
	exist := actual.(*limiterItem)

	atomic.StoreInt64(&exist.lastSeen, now)

	return exist.limiter
}

// AllowKey key级限流(快速失败)
// 适合高并发接口
func (s *Store) AllowKey(key string, r rate.Limit, burst int) bool {
	limiter := s.Get(key, r, burst)
	return limiter.Allow()
}

// WaitKey key级限流(平滑限流)
// 适合用户接口(体验更好)
func (s *Store) WaitKey(ctx context.Context, key string, r rate.Limit, burst int) error {
	limiter := s.Get(key, r, burst)
	return limiter.Wait(ctx)
}

// 后台定时清理过期key
// 避免: 内存泄漏,key无限增长
func (s *Store) clean() {
	// 退出时通知
	defer s.wg.Done()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {

		// 每分钟执行一次清理
		case <-ticker.C:
			expire := time.Now().Add(-s.ttl).Unix()

			s.m.Range(func(key, value any) bool {
				item := value.(*limiterItem)

				// 超过TTL→删除
				if atomic.LoadInt64(&item.lastSeen) < expire {
					s.m.Delete(key)
				}
				return true
			})

		// 服务关闭
		case <-s.stop:
			return
		}
	}
}
