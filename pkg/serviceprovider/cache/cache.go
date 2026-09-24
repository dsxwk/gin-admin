package cache

import (
	"context"
	"gin/common/ctxkey"
	"gin/common/flag"
	"gin/config"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/eventbus"
	"gin/pkg/serviceprovider/logger"
	"sync"
	"time"

	"github.com/goccy/go-json"
)

// Cache 缓存接口
type Cache interface {
	Set(key string, value any, expire time.Duration) error // 设置缓存
	Get(key string) (any, bool)                            // 获取缓存
	Delete(key string) error                               // 删除缓存
	Expire(key string) (any, time.Time, bool, error)       // 获取缓存过期时间
}

type CacheProxy struct {
	driver string
	c      Cache
	bus    *eventbus.Bus
	ctx    context.Context
	conf   *config.Config
}

func NewCacheProxy(driver string, c Cache, bus *eventbus.Bus, conf *config.Config) *CacheProxy {
	return &CacheProxy{
		driver: driver,
		c:      c,
		bus:    bus,
		conf:   conf,
	}
}

var (
	cacheInstance   = make(map[string]*CacheProxy)
	cacheInstanceMu sync.RWMutex
)

func NewCache(driver string, conf *config.Config) *CacheProxy {
	if driver == "" {
		driver = "memory"
	}

	cacheInstanceMu.Lock()
	defer cacheInstanceMu.Unlock()

	if instance, exists := cacheInstance[driver]; exists {
		return instance
	}

	var c *CacheProxy
	switch driver {
	case "redis":
		c = NewRedisCache(conf)
	case "disk":
		c = NewDiskCache(conf)
	case "memory":
		c = NewMemoryCache(conf)
	default:
		logger.NewLogger(conf).Fatal("不支持的缓存驱动: " + driver)
		return nil
	}

	cacheInstance[driver] = c
	flag.Infof("%s缓存初始化成功", driver)
	return c
}

func (p *CacheProxy) WithContext(ctx context.Context) *CacheProxy {
	proxy := &CacheProxy{
		driver: p.driver,
		c:      p.c,
		bus:    p.bus,
		ctx:    ctx,
		conf:   p.conf,
	}

	switch driver := p.c.(type) {
	case *RedisCache:
		proxy.c = driver.WithContext(ctx)
	case *MemoryCache:
		proxy.c = driver.WithContext(ctx)
	case *DiskCache:
		proxy.c = driver.WithContext(ctx)
	}

	return proxy
}

func (p *CacheProxy) Set(key string, value any, expire time.Duration) error {
	start := time.Now()
	err := p.c.Set(key, value, expire)
	p.publish("Set", key, value, time.Since(start))
	return err
}

func (p *CacheProxy) Get(key string) (any, bool) {
	start := time.Now()
	val, ok := p.c.Get(key)
	p.publish("Get", key, val, time.Since(start))
	return val, ok
}

func (p *CacheProxy) Delete(key string) error {
	start := time.Now()
	err := p.c.Delete(key)
	p.publish("Delete", key, nil, time.Since(start))
	return err
}

func (p *CacheProxy) Expire(key string) (any, time.Time, bool, error) {
	start := time.Now()
	val, exp, ok, err := p.c.Expire(key)
	p.publish("Expire", key, val, time.Since(start))
	return val, exp, ok, err
}

// Lock 获取带看门狗的Redis分布式锁
func (p *CacheProxy) Lock(ctx context.Context, key string, ttl time.Duration) (*LockResult, error) {
	if p == nil {
		return nil, ErrLockUnsupported
	}

	redisCache, ok := p.c.(*RedisCache)
	if !ok {
		return nil, ErrLockUnsupported
	}

	return redisCache.Lock(ctx, key, ttl)
}

// GetAs 获取缓存并转换为指定类型
func GetAs[T any](p *CacheProxy, key string) (T, bool) {
	var zero T
	if p == nil {
		return zero, false
	}

	value, ok := p.Get(key)
	if !ok {
		return zero, false
	}

	if result, ok := value.(T); ok {
		return result, true
	}

	data, err := encodeCacheValue(value)
	if err != nil {
		return zero, false
	}

	var result T
	if err = json.Unmarshal(data, &result); err != nil {
		return zero, false
	}

	return result, true
}

// Close 关闭底层缓存驱动
func (p *CacheProxy) Close() error {
	if p == nil || p.c == nil {
		return nil
	}

	closer, ok := p.c.(interface{ Close() error })
	if !ok {
		return nil
	}

	return closer.Close()
}

func (p *CacheProxy) publish(method, key string, val any, cost time.Duration) {
	if p.bus != nil && p.ctx != nil {
		traceId, ok := p.ctx.Value(ctxkey.TraceIDKey).(string)
		if !ok || traceId == "" {
			traceId = "unknown"
		}
		p.bus.Publish(p.ctx, debugger.TopicCache, debugger.CacheEvent{
			TraceID: traceId,
			Driver:  p.driver,
			Name:    method,
			Cmd:     key,
			Args:    val,
			Ms:      float64(cost.Nanoseconds()) / 1e6,
		})
	}
}

// Redis 获取Redis缓存实例
func Redis(conf *config.Config) *RedisCache {
	proxy := NewCache("redis", conf)
	if proxy == nil {
		return nil
	}

	redisCache, _ := proxy.c.(*RedisCache)
	return redisCache
}
