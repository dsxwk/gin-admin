package container

import (
	"gin/config"
	"gin/pkg/cache"
	"gin/pkg/db/connection"
	"gin/pkg/logger"
	"gorm.io/gorm"
	"sync"
)

type Container struct {
	Config      *config.Config
	Log         *logger.Logger
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
			Config:      config.NewConfig(),
			Log:         logger.NewLogger(),
			DB:          connection.Db{}.GetDB(),
			Cache:       cache.NewCache(),
			RedisCache:  cache.NewRedisCache(),
			MemoryCache: cache.NewMemoryCache(),
			DiskCache:   cache.NewDiskCache(),
		}
	})
	return instance
}

func GetContainer() *Container {
	return NewContainer()
}
