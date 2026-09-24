package ratelimit

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
)

// closeWait 关闭清理协程的最长等待时间
const closeWait = 3 * time.Second

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

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	// 清理协程必须保持非阻塞,超时兜底避免拖住优雅退出
	select {
	case <-done:
	case <-time.After(closeWait):
	}
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
	}
}

// allow key级快速失败限流
func (s *keyedStore) allow(key string, r rate.Limit, burst int) bool {
	return s.get(key, r, burst).Allow()
}

// wait key级平滑限流
func (s *keyedStore) wait(ctx context.Context, key string, r rate.Limit, burst int) error {
	err := s.get(key, r, burst).Wait(ctx)
	// 等待期间可能被清理,结束后刷新最后访问时间
	s.touch(key)

	return err
}

// touch 刷新key最后访问时间
func (s *keyedStore) touch(key string) {
	actual, ok := s.items.Load(key)
	if !ok {
		return
	}

	actual.(*limiterItem).lastSeen.Store(time.Now().Unix())
}

// clean 定时清理空闲key
func (s *keyedStore) clean() {
	defer s.wg.Done()

	interval := time.Minute
	if s.ttl < interval {
		// 清理周期取TTL的一半,避免空闲key最长存活接近2倍TTL
		interval = s.ttl / 2
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
