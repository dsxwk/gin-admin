package container

import (
	"fmt"
	"sync"
)

var defaultContainer = New()

// Default 获取默认服务容器
func Default() *Container {
	return defaultContainer
}

// Container 泛型服务容器
type Container struct {
	mu     sync.RWMutex
	values map[string]any
}

// New 创建服务容器
func New() *Container {
	return &Container{
		values: make(map[string]any),
	}
}

// Set 注册服务实例
func (c *Container) Set[T any](name string, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if current, exists := c.values[name]; exists {
		if _, ok := current.(T); !ok {
			panic(fmt.Sprintf("service %s type mismatch: %T != %T", name, current, value))
		}
	}

	c.values[name] = value
}

// Get 获取服务实例
func (c *Container) Get[T any](name string) T {
	var zero T

	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.values[name]
	if !exists {
		return zero
	}

	result, ok := value.(T)
	if !ok {
		panic(fmt.Sprintf("service %s type mismatch: %T != %T", name, value, zero))
	}
	return result
}

// Delete 删除服务实例
func (c *Container) Delete(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.values, name)
}
