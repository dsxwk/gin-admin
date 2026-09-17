package ratelimit

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

// limiterItem key对应的令牌桶
type limiterItem struct {
	limiter  *rate.Limiter // 令牌桶
	lastSeen atomic.Int64  // 最后访问时间
	r        rate.Limit    // 令牌生成速率
	burst    int           // 桶容量
}

// keyedStore key级令牌桶存储
type keyedStore struct {
	items     sync.Map       // key->limiterItem
	ttl       time.Duration  // 空闲多久后清理
	stop      chan struct{}  // 清理协程停止信号
	wg        sync.WaitGroup // 等待清理协程退出
	closeOnce sync.Once      // 确保只关闭一次
}

// newKeyedStore 创建key级令牌桶存储
func newKeyedStore(ttl time.Duration) *keyedStore {
	s := &keyedStore{
		ttl:  ttl,
		stop: make(chan struct{}),
	}
	if ttl > 0 {
		s.wg.Add(1)
		go s.clean()
	}
	return s
}

// close 关闭清理协程
func (s *keyedStore) close() {
	if s == nil {
		return
	}
	s.closeOnce.Do(func() {
		if s.stop != nil {
			close(s.stop)
		}
	})
	s.wg.Wait()
}

// get 获取key对应的令牌桶
func (s *keyedStore) get(key string, r rate.Limit, burst int) *rate.Limiter {
	now := time.Now().Unix()
	candidate := newLimiterItem(r, burst, now)

	for {
		actual, loaded := s.items.LoadOrStore(key, candidate)
		if !loaded {
			return candidate.limiter
		}

		current := actual.(*limiterItem)
		if current.r == r && current.burst == burst {
			current.lastSeen.Store(now)
			return current.limiter
		}

		replacement := newLimiterItem(r, burst, now)
		if s.items.CompareAndSwap(key, current, replacement) {
			return replacement.limiter
		}
		candidate = newLimiterItem(r, burst, now)
	}
}

// allow key级快速失败限流
func (s *keyedStore) allow(key string, r rate.Limit, burst int) bool {
	return s.get(key, r, burst).Allow()
}

// wait key级平滑限流
func (s *keyedStore) wait(ctx context.Context, key string, r rate.Limit, burst int) error {
	return s.get(key, r, burst).Wait(ctx)
}

// clean 定时清理空闲key
func (s *keyedStore) clean() {
	defer s.wg.Done()

	interval := time.Minute
	if s.ttl < interval {
		interval = s.ttl
	}
	if interval < time.Millisecond {
		interval = time.Millisecond
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			expire := time.Now().Add(-s.ttl).Unix()
			s.items.Range(func(key, value any) bool {
				item := value.(*limiterItem)
				if item.lastSeen.Load() < expire {
					s.items.Delete(key)
				}
				return true
			})
		case <-s.stop:
			return
		}
	}
}

// newLimiterItem 创建key对应的令牌桶
func newLimiterItem(r rate.Limit, burst int, lastSeen int64) *limiterItem {
	item := &limiterItem{
		limiter: rate.NewLimiter(r, burst),
		r:       r,
		burst:   burst,
	}
	item.lastSeen.Store(lastSeen)
	return item
}
