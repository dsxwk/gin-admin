package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"gin/common/flag"
	"gin/config"
	"os"
	"sync"
)

// Consumer interface for queue consumers
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

// ConsumerStatus consumer status
type ConsumerStatus string

const (
	ConsumerStatusStopped ConsumerStatus = "stopped"
	ConsumerStatusRunning ConsumerStatus = "running"
	ConsumerStatusError   ConsumerStatus = "error"
)

// PayloadHandler payload handler interface
type PayloadHandler interface {
	NewPayload() any
	Handle(payload any) error
}

// Producer interface for queue producers
type Producer interface {
	Name() string
	Description() string
	Connection() string
	IsDelay() bool
	DelayMs() int64
	Publish(ctx context.Context, msg any) error
	Close() error
}

// Registry generic registry
type Registry[T any] struct {
	items map[string]T
	mu    sync.RWMutex
}

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{items: make(map[string]T)}
}

type namer interface {
	Name() string
}

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

// TryHandle auto deserialize and call Handle
func TryHandle(h interface{}, body []byte) error {
	ph, ok := h.(PayloadHandler)
	if !ok {
		return fmt.Errorf("consumer does not implement PayloadHandler")
	}
	payload := ph.NewPayload()
	if err := json.Unmarshal(body, payload); err != nil {
		return err
	}
	return ph.Handle(payload)
}
