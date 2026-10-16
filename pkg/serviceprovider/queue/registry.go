package queue

import (
	"fmt"
	"sync"
)

// Registry 泛型注册表
type Registry[T Named] struct {
	items map[string]T
	mu    sync.RWMutex
}

// NewRegistry 创建泛型注册表
func NewRegistry[T Named]() *Registry[T] {
	return &Registry[T]{items: make(map[string]T)}
}

// Register 注册项目
func (r *Registry[T]) Register(item T) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := item.Name()
	if _, exists := r.items[name]; exists {
		return fmt.Errorf("queue %s already registered", name)
	}

	r.items[name] = item
	return nil
}

// Get 获取项目
func (r *Registry[T]) Get(name string) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[name]
	return item, exists
}

// All 获取全部项目
func (r *Registry[T]) All() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]T, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	return items
}

// Names 获取全部名称
func (r *Registry[T]) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.items))
	for name := range r.items {
		names = append(names, name)
	}
	return names
}

// Count 获取项目数量
func (r *Registry[T]) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.items)
}

// Exists 判断项目是否存在
func (r *Registry[T]) Exists(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.items[name]
	return exists
}
