package container

import (
	"gin/config"
	"gin/pkg/cache"
	"gorm.io/gorm"
	"sync"
)

type Container struct {
	Config      *config.Config
	Log         *config.Logger
	DB          *gorm.DB
	Cache       *cache.CacheProxy
	RedisCache  *cache.CacheProxy
	MemoryCache *cache.CacheProxy
	DiskCache   *cache.CacheProxy
}

var (
	instance *Container
	once     sync.Once
)

func NewContainer() *Container {
	once.Do(func() {
		instance = &Container{
			Config:      config.GetConfig(),
			Log:         config.GetLogger(),
			DB:          config.Db{}.GetDB(),
			Cache:       config.GetCache(),
			RedisCache:  config.GetRedisCache(),
			MemoryCache: config.GetMemoryCache(),
			DiskCache:   config.GetDiskCache(),
		}
	})
	return instance
}

func GetContainer() *Container {
	return NewContainer()
}
