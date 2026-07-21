package cache

import (
	"context"
	"gin/common/ctxkey"
	"gin/common/flag"
	"gin/config"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/logger"
	"gin/pkg/serviceprovider/message"
	"sync"
	"time"
)

// Cache 缓存接口
type Cache interface {
	Set(key string, value interface{}, expire time.Duration) error // 设置缓存
	Get(key string) (interface{}, bool)                            // 获取缓存
	Delete(key string) error                                       // 删除缓存
	Expire(key string) (interface{}, time.Time, bool, error)       // 获取缓存过期时间
}

type CacheProxy struct {
	driver string
	c      Cache
	bus    *message.Event
	ctx    context.Context
	conf   *config.Config
}

func NewCacheProxy(driver string, c Cache, bus *message.Event, conf *config.Config) *CacheProxy {
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
	cacheInstanceMu.RLock()
	c, ok := cacheInstance[driver]
	cacheInstanceMu.RUnlock()
	if ok {
		return c
	}
	switch driver {
	case "redis", "disk":
		c = NewRedisCache(conf)
		cacheInstanceMu.Lock()
		cacheInstance[driver] = c
		cacheInstanceMu.Unlock()
		flag.Infof("%s缓存初始化成功", driver)
		return c
	case "", "memory":
		c = NewMemoryCache(conf)
		cacheInstanceMu.Lock()
		cacheInstance["memory"] = c
		cacheInstanceMu.Unlock()
		flag.Infof("memory缓存初始化成功")
		return c
	default:
		logger.NewLogger(conf).Fatal("不支持的缓存驱动: " + driver)
		return nil
	}
}

func ResetCache(driver string, conf *config.Config) *CacheProxy {
	cacheInstanceMu.Lock()
	delete(cacheInstance, driver)
	delete(cacheInstance, conf.Cache.Driver)
	cacheInstanceMu.Unlock()
	return NewCache(driver, conf)
}
func (p *CacheProxy) WithContext(ctx context.Context) *CacheProxy {
	return &CacheProxy{
		driver: p.driver,
		c:      p.c,
		bus:    p.bus,
		ctx:    ctx,
	}
}

func (p *CacheProxy) Set(key string, value interface{}, expire time.Duration) error {
	start := time.Now()
	err := p.c.Set(key, value, expire)
	p.publish("Set", key, value, time.Since(start))
	return err
}

func (p *CacheProxy) Get(key string) (interface{}, bool) {
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

func (p *CacheProxy) Expire(key string) (interface{}, time.Time, bool, error) {
	start := time.Now()
	val, exp, ok, err := p.c.Expire(key)
	p.publish("Expire", key, val, time.Since(start))
	return val, exp, ok, err
}

func (p *CacheProxy) publish(method, key string, val interface{}, cost time.Duration) {
	if p.bus != nil && p.ctx != nil {
		traceId, ok := p.ctx.Value(ctxkey.TraceIdKey).(string)
		if !ok || traceId == "" {
			traceId = "unknown"
		}
		p.bus.Publish(debugger.TopicCache, debugger.CacheEvent{
			TraceId: traceId,
			Driver:  p.driver,
			Name:    method,
			Cmd:     key,
			Args:    val,
			Ms:      float64(cost.Nanoseconds()) / 1e6,
		})
	}
}

func (p *CacheProxy) Redis() *RedisCache {
	r := p.c.(*RedisCache)
	if err := r.Ping(); err != nil {
		ResetRedisCache(p.conf)
		p.c = redisCache.c
		r = p.c.(*RedisCache)
	}
	return r
}
