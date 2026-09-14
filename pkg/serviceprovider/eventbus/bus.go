package eventbus

import (
	"context"
	"sync"
	"sync/atomic"
)

// Handler 泛型事件处理函数
type Handler[T any] func(T)

type busSubscriber struct {
	id     uint64
	handle func(any)
	async  bool
}

// Bus 底层发布订阅总线
type Bus struct {
	mu          sync.RWMutex
	subscribers map[string][]*busSubscriber
	nextID      atomic.Uint64
}

var (
	defaultBus *Bus
	busOnce    sync.Once
)

// NewBus 获取事件总线单例
func NewBus() *Bus {
	busOnce.Do(func() {
		defaultBus = &Bus{
			subscribers: make(map[string][]*busSubscriber),
		}
	})
	return defaultBus
}

// Subscribe 订阅同步事件
func (b *Bus) Subscribe[T any](topic string, handler Handler[T]) uint64 {
	return b.subscribe(topic, handler, false)
}

// SubscribeAsync 订阅异步事件
func (b *Bus) SubscribeAsync[T any](topic string, handler Handler[T]) uint64 {
	return b.subscribe(topic, handler, true)
}

// subscribe 创建指定模式的事件订阅
func (b *Bus) subscribe[T any](topic string, handler Handler[T], async bool) uint64 {
	if b == nil || handler == nil {
		return 0
	}

	return b.addSubscriber(topic, func(data any) {
		event, ok := data.(T)
		if !ok {
			return
		}

		handler(event)
	}, async)
}

// addSubscriber 添加订阅者
func (b *Bus) addSubscriber(topic string, handler func(any), async bool) uint64 {
	id := b.nextID.Add(1)
	subscriber := &busSubscriber{
		id:     id,
		handle: handler,
		async:  async,
	}

	b.mu.Lock()
	if b.subscribers == nil {
		b.subscribers = make(map[string][]*busSubscriber)
	}
	b.subscribers[topic] = append(b.subscribers[topic], subscriber)
	b.mu.Unlock()

	return id
}

// Unsubscribe 取消指定订阅
func (b *Bus) Unsubscribe(topic string, id uint64) bool {
	if b == nil {
		return false
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	subscribers := b.subscribers[topic]
	if len(subscribers) == 0 {
		return false
	}

	for index, subscriber := range subscribers {
		if subscriber.id != id {
			continue
		}

		b.subscribers[topic] = append(subscribers[:index], subscribers[index+1:]...)
		if len(b.subscribers[topic]) == 0 {
			delete(b.subscribers, topic)
		}

		return true
	}

	return false
}

// Publish 发布事件
func (b *Bus) Publish(topic string, event any) {
	b.PublishWithContext(context.Background(), topic, event)
}

// PublishWithContext 发布事件上下文
func (b *Bus) PublishWithContext(ctx context.Context, topic string, event any) {
	if b == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}

	subscribers := b.subscribersFor(topic)
	for _, subscriber := range subscribers {
		if ctx.Err() != nil {
			return
		}

		if subscriber.async {
			b.asyncHandle(ctx, subscriber, event)
			continue
		}

		b.handle(ctx, subscriber, event)
	}
}

// subscribersFor 获取订阅者快照
func (b *Bus) subscribersFor(topic string) []*busSubscriber {
	if b == nil {
		return nil
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	subscribers := b.subscribers[topic]
	if len(subscribers) == 0 {
		return nil
	}

	result := make([]*busSubscriber, len(subscribers))
	copy(result, subscribers)
	return result
}

// handle 执行订阅者
func (b *Bus) handle(ctx context.Context, subscriber *busSubscriber, event any) {
	defer func() {
		_ = recover()
	}()

	if ctx.Err() != nil {
		return
	}

	subscriber.handle(event)
}

// asyncHandle 异步执行订阅者
func (b *Bus) asyncHandle(ctx context.Context, subscriber *busSubscriber, event any) {
	go func() {
		b.handle(ctx, subscriber, event)
	}()
}

// Count 获取指定主题的订阅数量
func (b *Bus) Count(topic string) int {
	if b == nil {
		return 0
	}

	b.mu.RLock()
	defer b.mu.RUnlock()

	return len(b.subscribers[topic])
}
