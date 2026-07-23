package queue

import (
	"context"
	"gin/common/flag"
	"gin/config"
	"os"
	"sync"
)

// Consumer 消费者接口
// 消费者在应用启动时自动启动
type Consumer interface {
	// Name 消费者名称
	Name() string
	// Description 消费者描述
	Description() string
	// Start 启动消费者
	Start() error
	// Stop 停止消费者
	Stop() error
	// Enabled 是否启用
	Enabled(cfg *config.Config) bool
	// Status 消费者状态
	Status() ConsumerStatus
}

// ConsumerStatus 消费者状态
type ConsumerStatus string

const (
	ConsumerStatusStopped ConsumerStatus = "stopped" // 已停止
	ConsumerStatusRunning ConsumerStatus = "running" // 运行中
	ConsumerStatusError   ConsumerStatus = "error"   // 错误
)

// Handler 消息处理接口
// 消费者需要实现此接口来处理业务逻辑
type Handler interface {
	Handle(msg string) error
}

// Producer 生产者接口
// 生产者按需使用,不需要在应用启动时初始化
type Producer interface {
	// Name 生产者名称
	Name() string
	// Description 生产者描述
	Description() string
	// Publish 发送消息
	Publish(ctx context.Context, msg []byte) error
	// Close 关闭生产者
	Close() error
}

// Registry 通用注册表
type Registry[T any] struct {
	items map[string]T
	mu    sync.RWMutex
}

// NewRegistry 创建注册表
func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{
		items: make(map[string]T),
	}
}

// namer 名称提取接口
type namer interface {
	Name() string
}

// Register 注册
func (r *Registry[T]) Register(item T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := ""
	if n, ok := any(item).(namer); ok {
		name = n.Name()
	}

	if name == "" {
		flag.Errorf("Queue name cannot be empty")
		os.Exit(1)
	}

	if _, exists := r.items[name]; exists {
		flag.Errorf("Queue %s already registered", name)
		os.Exit(1)
	}

	r.items[name] = item
}

// Get 获取指定名称的项
func (r *Registry[T]) Get(name string) T {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.items[name]
}

// GetAll 获取所有已注册的项
func (r *Registry[T]) GetAll() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]T, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	return items
}

// GetNames 获取所有名称列表
func (r *Registry[T]) GetNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.items))
	for name := range r.items {
		names = append(names, name)
	}
	return names
}

// Count 获取数量
func (r *Registry[T]) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.items)
}

// Exists 检查是否存在
func (r *Registry[T]) Exists(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.items[name]
	return exists
}

var (
	// Consumers 消费者注册表
	Consumers = NewRegistry[Consumer]()
	// Producers 生产者注册表
	Producers = NewRegistry[Producer]()
)

// GetConsumerRegistry 获取消费者注册表
func GetConsumerRegistry() *Registry[Consumer] {
	return Consumers
}

// GetProducerRegistry 获取生产者注册表
func GetProducerRegistry() *Registry[Producer] {
	return Producers
}
