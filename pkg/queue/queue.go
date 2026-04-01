package queue

import (
	"context"
	"gin/common/flag"
	"gin/config"
	"gin/pkg/logger"
	"os"
	"sync"
)

// Consumer 消费者接口
// 消费者在应用启动时自动启动
type Consumer interface {
	// Name 消费者名称
	Name() string
	// Start 启动消费者
	Start(cfg *config.Config, log *logger.Logger) error
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
	// Publish 发送消息
	Publish(ctx context.Context, msg []byte) error
	// Close 关闭生产者
	Close() error
}

// ConsumerRegistry 消费者注册表
// 用于管理所有已注册的消费者
type ConsumerRegistry struct {
	consumers map[string]Consumer
	mu        sync.RWMutex
}

var consumerRegistry = &ConsumerRegistry{
	consumers: make(map[string]Consumer),
}

// Register 注册消费者(在消费者的init函数中调用)
func (r *ConsumerRegistry) Register(c Consumer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := c.Name()
	if name == "" {
		flag.Errorf("Consumer name cannot be empty")
		os.Exit(1)
	}

	// 检查是否已存在
	if _, exists := r.consumers[name]; exists {
		flag.Errorf("Consumer %s already registered", name)
		os.Exit(1)
	}

	r.consumers[name] = c
}

// Get 获取指定名称的消费者
func (r *ConsumerRegistry) Get(name string) Consumer {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.consumers[name]
}

// GetAll 获取所有已注册的消费者
func (r *ConsumerRegistry) GetAll() []Consumer {
	r.mu.RLock()
	defer r.mu.RUnlock()

	consumers := make([]Consumer, 0, len(r.consumers))
	for _, c := range r.consumers {
		consumers = append(consumers, c)
	}
	return consumers
}

// GetNames 获取所有消费者名称列表
func (r *ConsumerRegistry) GetNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.consumers))
	for name := range r.consumers {
		names = append(names, name)
	}
	return names
}

// Count 获取消费者数量
func (r *ConsumerRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.consumers)
}

// Exists 检查消费者是否存在
func (r *ConsumerRegistry) Exists(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.consumers[name]
	return exists
}

// ProducerRegistry 生产者注册表
// 用于管理所有已注册的生产者
type ProducerRegistry struct {
	producers map[string]Producer
	mu        sync.RWMutex
}

var producerRegistry = &ProducerRegistry{
	producers: make(map[string]Producer),
}

// Register 注册生产者(在生产者的init函数中调用)
func (r *ProducerRegistry) Register(p Producer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := p.Name()
	if name == "" {
		flag.Errorf("Producer name cannot be empty")
		os.Exit(1)
	}

	// 检查是否已存在
	if _, exists := r.producers[name]; exists {
		flag.Errorf("Producer %s already registered", name)
		os.Exit(1)
	}

	r.producers[name] = p
}

// Get 获取指定名称的生产者
func (r *ProducerRegistry) Get(name string) Producer {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.producers[name]
}

// GetAll 获取所有已注册的生产者
func (r *ProducerRegistry) GetAll() []Producer {
	r.mu.RLock()
	defer r.mu.RUnlock()

	producers := make([]Producer, 0, len(r.producers))
	for _, p := range r.producers {
		producers = append(producers, p)
	}
	return producers
}

// GetNames 获取所有生产者名称列表
func (r *ProducerRegistry) GetNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.producers))
	for name := range r.producers {
		names = append(names, name)
	}
	return names
}

// Count 获取生产者数量
func (r *ProducerRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.producers)
}

// Exists 检查生产者是否存在
func (r *ProducerRegistry) Exists(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.producers[name]
	return exists
}

// GetConsumerRegistry 获取消费者注册表
func GetConsumerRegistry() *ConsumerRegistry {
	return consumerRegistry
}

// GetProducerRegistry 获取生产者注册表
func GetProducerRegistry() *ProducerRegistry {
	return producerRegistry
}
