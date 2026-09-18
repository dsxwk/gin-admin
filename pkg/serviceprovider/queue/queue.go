package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gin/config"
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
	ConsumerStatusStopped  ConsumerStatus = "stopped"
	ConsumerStatusStarting ConsumerStatus = "starting"
	ConsumerStatusRunning  ConsumerStatus = "running"
	ConsumerStatusError    ConsumerStatus = "error"
)

// PayloadHandler 消息负载处理接口
type PayloadHandler interface {
	NewPayload() any
	Handle(payload any) error
}

// ContextPayloadHandler 支持上下文的消息处理接口
type ContextPayloadHandler interface {
	PayloadHandler
	HandleContext(ctx context.Context, payload any) error
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
	items     map[string]T
	factories []func() T
	mu        sync.RWMutex
}

func NewRegistry[T Named]() *Registry[T] {
	return &Registry[T]{items: make(map[string]T)}
}

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

// RegisterFactory 注册延迟创建工厂
func (r *Registry[T]) RegisterFactory(factory func() T) {
	if factory == nil {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.factories = append(r.factories, factory)
}

// RegisterFactories 创建并注册全部延迟实例
func (r *Registry[T]) RegisterFactories() error {
	r.mu.Lock()
	factories := r.factories
	r.factories = nil
	r.mu.Unlock()

	var errs []error
	for _, factory := range factories {
		item := factory()
		if any(item) == nil {
			continue
		}
		if err := r.Register(item); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (r *Registry[T]) Get(name string) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[name]
	return item, exists
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

// TryHandle 自动反序列化并调用处理,支持上下文
func TryHandle[T PayloadHandler](ctx context.Context, h T, body []byte) error {
	if ctx == nil {
		ctx = context.Background()
	}
	payload := h.NewPayload()
	if err := json.Unmarshal(body, payload); err != nil {
		return err
	}
	if handler, ok := any(h).(ContextPayloadHandler); ok {
		return handler.HandleContext(ctx, payload)
	}
	return h.Handle(payload)
}
