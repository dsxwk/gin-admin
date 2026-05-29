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
}

func NewCacheProxy(driver string, c Cache, bus *message.Event) *CacheProxy {
	return &CacheProxy{
		driver: driver,
		c:      c,
		bus:    bus,
	}
}

var (
	cacheInstance = make(map[string]*CacheProxy)
	cacheLocks    sync.Map
)

func NewCache(driver string, conf *config.Config) *CacheProxy {
	// 每个driver只初始化一次
	onceAny, _ := cacheLocks.LoadOrStore(driver, &sync.Once{})
	once := onceAny.(*sync.Once)
	_default := "memory"
	once.Do(func() {
		switch driver {

		case "redis", "disk":
			_default = driver
			cacheInstance[_default] = NewRedisCache(conf)

		case "", "memory":
			_default = "memory"
			cacheInstance["_default"] = NewMemoryCache(conf)

		default:
			logger.NewLogger(conf).Fatal("不支持的缓存驱动: " + driver)
		}

		flag.Infof("%s缓存初始化成功", _default)
	})

	return cacheInstance[_default]
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
	// 强制类型断言
	return p.c.(*RedisCache)
}
