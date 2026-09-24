package cache

import (
	"context"
	"fmt"
	"gin/pkg/errcode"
	"sync"
	"time"
	"uuid"

	"github.com/go-redis/redis/v8"
)

const (
	// CacheErrCodePrefix 缓存错误码前缀
	CacheErrCodePrefix = 700
	lockReleaseScript  = `if redis.call("get",KEYS[1])==ARGV[1] then
	return redis.call("del",KEYS[1])
else
	return 0
end`
	lockRenewScript = `if redis.call("get",KEYS[1])==ARGV[1] then
	return redis.call("pexpire",KEYS[1],ARGV[2])
else
	return 0
end`
	lockOperationTimeout = 5 * time.Second
	lockMinTTL           = 10 * time.Millisecond
)

var (
	// ErrLockExists 锁已存在
	ErrLockExists = errcode.NewError(1, "锁已存在").WithPrefix(CacheErrCodePrefix)
	// ErrLockNotOwned 锁不属于当前持有者或已释放
	ErrLockNotOwned = errcode.NewError(2, "锁不属于当前持有者或已释放").WithPrefix(CacheErrCodePrefix)
	// ErrLockInvalidTTL 锁的过期时间不能小于10毫秒
	ErrLockInvalidTTL = errcode.NewError(3, "锁的过期时间不能小于10毫秒").WithPrefix(CacheErrCodePrefix)
	// ErrLockUnsupported 当前缓存驱动不支持分布式锁
	ErrLockUnsupported = errcode.NewError(4, "当前缓存驱动不支持分布式锁").WithPrefix(CacheErrCodePrefix)
	// ErrCacheClosed 缓存服务已关闭
	ErrCacheClosed = errcode.NewError(5, "缓存服务已关闭").WithPrefix(CacheErrCodePrefix)
	// ErrLockAcquire 获取锁失败
	ErrLockAcquire = errcode.NewError(6, "获取锁失败").WithPrefix(CacheErrCodePrefix)
	// ErrLockRelease 释放锁失败
	ErrLockRelease = errcode.NewError(7, "释放锁失败").WithPrefix(CacheErrCodePrefix)
)

// LockResult Redis锁持有结果
type LockResult struct {
	client *redis.Client
	state  *redisState
	key    string
	token  string
	ttl    time.Duration
	stop   chan struct{}
	done   chan struct{}
	once   sync.Once
	err    error
}

// Lock 获取锁并启动看门狗续期
func (r *RedisCache) Lock(ctx context.Context, key string, ttl time.Duration) (*LockResult, error) {
	if r == nil || r.client == nil {
		return nil, ErrLockNotOwned
	}
	if ttl < lockMinTTL {
		return nil, ErrLockInvalidTTL
	}
	if ctx == nil {
		ctx = context.Background()
	}

	token := uuid.New().String()
	ok, err := r.client.SetNX(ctx, key, token, ttl).Result()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLockAcquire, err)
	}
	if !ok {
		return nil, ErrLockExists
	}

	result := &LockResult{
		client: r.client,
		state:  r.state,
		key:    key,
		token:  token,
		ttl:    ttl,
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}
	if r.state != nil {
		r.state.mu.Lock()
		if r.state.closed {
			r.state.mu.Unlock()
			_ = r.releaseLockToken(context.Background(), key, token)
			return nil, ErrCacheClosed
		}
		r.state.locks[key] = result
		r.state.mu.Unlock()
	}
	go result.watchdog()

	return result, nil
}

// releaseLockToken 校验token并释放锁
func (r *RedisCache) releaseLockToken(ctx context.Context, key string, token string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx = context.WithoutCancel(ctx)
	ctx, cancel := context.WithTimeout(ctx, lockOperationTimeout)
	defer cancel()

	status, err := r.client.Eval(ctx, lockReleaseScript, []string{key}, token).Int()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrLockRelease, err)
	}
	if status == 0 {
		return ErrLockNotOwned
	}

	return nil
}

// watchdog 自动续期
func (l *LockResult) watchdog() {
	defer close(l.done)

	interval := max(l.ttl/3, time.Millisecond)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-l.stop:
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), lockOperationTimeout)
			status, err := l.client.Eval(
				ctx,
				lockRenewScript,
				[]string{l.key},
				l.token,
				l.ttl.Milliseconds(),
			).Int()
			cancel()

			if err == nil && status == 0 {
				return
			}
		}
	}
}

// Release 停止续期并原子释放锁
func (l *LockResult) Release() error {
	if l == nil || l.client == nil {
		return ErrLockNotOwned
	}

	l.once.Do(func() {
		defer l.unregister()

		close(l.stop)
		<-l.done

		ctx, cancel := context.WithTimeout(context.Background(), lockOperationTimeout)
		defer cancel()

		status, err := l.client.Eval(ctx, lockReleaseScript, []string{l.key}, l.token).Int()
		if err != nil {
			l.err = fmt.Errorf("%w: %v", ErrLockRelease, err)
			return
		}
		if status == 0 {
			l.err = ErrLockNotOwned
		}
	})

	return l.err
}

// unregister 移除锁持有记录
func (l *LockResult) unregister() {
	if l == nil || l.state == nil {
		return
	}

	l.state.mu.Lock()
	if current := l.state.locks[l.key]; current == l {
		delete(l.state.locks, l.key)
	}
	l.state.mu.Unlock()
}

// Key 获取锁的key
func (l *LockResult) Key() string {
	if l == nil {
		return ""
	}

	return l.key
}
