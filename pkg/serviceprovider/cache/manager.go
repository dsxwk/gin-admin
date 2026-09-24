package cache

import (
	"gin/config"
	"sync"
)

// Manager 缓存管理器
type Manager struct {
	mu          sync.RWMutex
	conf        *config.Config
	defaultName string
	caches      map[string]*CacheProxy
}

// NewManager 创建缓存管理器
func NewManager(conf *config.Config) *Manager {
	name := "memory"
	if conf != nil && conf.Cache.Driver != "" {
		name = conf.Cache.Driver
	}

	return &Manager{
		conf:        conf,
		defaultName: name,
		caches:      make(map[string]*CacheProxy),
	}
}

// Cache 获取指定缓存实例
func (m *Manager) Cache(name ...string) *CacheProxy {
	driver := m.defaultName
	if len(name) > 0 && name[0] != "" {
		driver = name[0]
	}

	m.mu.RLock()
	instance := m.caches[driver]
	m.mu.RUnlock()
	if instance != nil {
		return instance
	}

	instance = NewCache(driver, m.conf)
	m.mu.Lock()
	m.caches[driver] = instance
	m.mu.Unlock()
	return instance
}

// Redis 获取Redis缓存
func (m *Manager) Redis() *RedisCache {
	instance := m.Cache("redis")
	if instance == nil {
		return nil
	}

	redisCache, _ := instance.c.(*RedisCache)
	return redisCache
}

// Close 关闭全部缓存实例
func (m *Manager) Close() error {
	if m == nil {
		return nil
	}

	m.mu.Lock()
	instances := make([]*CacheProxy, 0, len(m.caches))
	for _, instance := range m.caches {
		instances = append(instances, instance)
	}
	m.caches = make(map[string]*CacheProxy)
	m.mu.Unlock()

	var firstErr error
	for _, instance := range instances {
		if err := instance.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}
