package queue

import (
	"context"
	"encoding/json"
	"gin/common/flag"
	"gin/config"
	"os"
	"sync"
)

// Consumer 队列消费者接口
type Consumer interface {
	Name() string
	Description() string
	Connection() string
	Retry() int
	IsDelay() bool
	Start() error
	Stop() error
	Enabled(cfg *config.Config) bool
	Status() ConsumerStatus
}

// ConsumerStatus 消费者状态
type ConsumerStatus string

const (
	ConsumerStatusStopped ConsumerStatus = "stopped"
	ConsumerStatusRunning ConsumerStatus = "running"
	ConsumerStatusError   ConsumerStatus = "error"
)

// PayloadHandler 消息负载处理接口
type PayloadHandler interface {
	NewPayload() any
	Handle(payload any) error
}

// Producer 队列生产者接口
type Producer interface {
	Name() string
	Description() string
	Connection() string
	IsDelay() bool
	DelayMs() int64
	Publish(ctx context.Context, msg any) error
	Close() error
}

// Named 名称接口
type Named interface {
	Name() string
}

// ConsumerHandler 消费者处理接口
type ConsumerHandler interface {
	Consumer
	PayloadHandler
}

// Registry 泛型注册表
type Registry[T Named] struct {
	items map[string]T
	mu    sync.RWMutex
}

func NewRegistry[T Named]() *Registry[T] {
	return &Registry[T]{items: make(map[string]T)}
}

func (r *Registry[T]) Register(item T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := item.Name()

	if _, exists := r.items[name]; exists {
		flag.Errorf("Queue %s already registered", name)
		os.Exit(1)
	}

	r.items[name] = item
}

func (r *Registry[T]) Get(name string) T {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.items[name]
}

func (r *Registry[T]) GetAll() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]T, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}
	return items
}

func (r *Registry[T]) GetNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.items))
	for name := range r.items {
		names = append(names, name)
	}
	return names
}

func (r *Registry[T]) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.items)
}

func (r *Registry[T]) Exists(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.items[name]
	return exists
}

var (
	Consumers = NewRegistry[Consumer]()
	Producers = NewRegistry[Producer]()
)

func GetConsumerRegistry() *Registry[Consumer] {
	return Consumers
}

func GetProducerRegistry() *Registry[Producer] {
	return Producers
}

// TryHandle 自动反序列化并调用处理
func TryHandle[T PayloadHandler](h T, body []byte) error {
	payload := h.NewPayload()
	if err := json.Unmarshal(body, payload); err != nil {
		return err
	}
	return h.Handle(payload)
}
