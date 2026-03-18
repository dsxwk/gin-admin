package container

import (
	"context"
	"gin/common/ctxkey"
)

func (c *Container) WithContext(ctx context.Context) *Container {
	return &Container{
		Config:      c.Config,
		Log:         c.Log,
		DB:          c.DB.WithContext(ctx),
		Cache:       c.Cache.WithContext(ctx),
		RedisCache:  c.RedisCache.WithContext(ctx),
		MemoryCache: c.MemoryCache.WithContext(ctx),
		DiskCache:   c.DiskCache.WithContext(ctx),
	}
}

// Set 保存Container到Context
func Set(ctx context.Context, c *Container) context.Context {
	return context.WithValue(ctx, ctxkey.ContainerKey, c)
}

// Get 从Context获取Container
func Get(ctx context.Context) *Container {
	if c, ok := ctx.Value(ctxkey.ContainerKey).(*Container); ok {
		return c
	}
	return GetContainer()
}
