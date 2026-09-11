package eventbus

import (
	"context"
	"sync"
	"sync/atomic"
)

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

// busSubscriber 总线订阅者
type busSubscriber struct {
	Id     uint64
	Async  bool
	Handle func(any)
}

// Bus 底层事件总线
type Bus struct {
	subscribers map[string][]*busSubscriber
	mu          sync.RWMutex
	idCounter   atomic.Uint64
	semaphore   chan struct{}
}

// Subscribe 订阅同步事件
func (b *Bus) Subscribe[T any](topic string, fn func(T)) uint64 {
	return b.addSubscriber(topic, func(event any) {
		if v, ok := event.(T); ok {
			fn(v)
		}
	}, false)
}

// SubscribeAsync 订阅异步事件
func (b *Bus) SubscribeAsync[T any](topic string, fn func(T)) uint64 {
	return b.addSubscriber(topic, func(event any) {
		if v, ok := event.(T); ok {
			fn(v)
		}
	}, true)
}

// addSubscriber 添加订阅者
func (b *Bus) addSubscriber(topic string, fn func(any), async bool) uint64 {
	id := b.idCounter.Add(1)

	sub := &busSubscriber{
		Id:     id,
		Async:  async,
		Handle: fn,
	}

	b.mu.Lock()
	b.subscribers[topic] = append(b.subscribers[topic], sub)
	b.mu.Unlock()

	return id
}

// Unsubscribe 删除订阅者
func (b *Bus) Unsubscribe(topic string, id uint64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	subs, ok := b.subscribers[topic]
	if !ok {
		return false
	}

	newSubs := make([]*busSubscriber, 0, len(subs))
	removed := false

	for _, s := range subs {
		if s.Id == id {
			removed = true
			continue
		}
		newSubs = append(newSubs, s)
	}

	if removed {
		if len(newSubs) == 0 {
			delete(b.subscribers, topic)
		} else {
			b.subscribers[topic] = newSubs
		}
	}

	return removed
}

// UnsubscribeAll 取消主题的所有订阅
func (b *Bus) UnsubscribeAll(topic string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, ok := b.subscribers[topic]; !ok {
		return false
	}

	delete(b.subscribers, topic)
	return true
}

// Publish 发布事件
func (b *Bus) Publish[T any](topic string, event T) {
	b.mu.RLock()
	subs := b.subscribers[topic]
	subsCopy := make([]*busSubscriber, len(subs))
	copy(subsCopy, subs)
	b.mu.RUnlock()

	if len(subsCopy) == 0 {
		return
	}

	for _, sub := range subsCopy {
		if sub.Async {
			b.asyncHandle(sub, event)
		} else {
			sub.Handle(event)
		}
	}
}

// PublishWithContext 发布事件并支持取消
func (b *Bus) PublishWithContext[T any](ctx context.Context, topic string, event T) {
	b.mu.RLock()
	subs := b.subscribers[topic]
	subsCopy := make([]*busSubscriber, len(subs))
	copy(subsCopy, subs)
	b.mu.RUnlock()

	if len(subsCopy) == 0 {
		return
	}

	for _, sub := range subsCopy {
		select {
		case <-ctx.Done():
			return
		default:
			if sub.Async {
				b.asyncHandle(sub, event)
			} else {
				sub.Handle(event)
			}
		}
	}
}

// asyncHandle 异步执行
func (b *Bus) asyncHandle(sub *busSubscriber, event any) {
	if b.semaphore != nil {
		b.semaphore <- struct{}{}
		go func() {
			defer func() { <-b.semaphore }()
			sub.Handle(event)
		}()
		return
	}
	sub.Handle(event)
}

// SubscribeIds 查询主题的所有订阅者ID
func (b *Bus) SubscribeIds(topic string) []uint64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	subs, ok := b.subscribers[topic]
	if !ok {
		return []uint64{}
	}

	ids := make([]uint64, 0, len(subs))
	for _, s := range subs {
		ids = append(ids, s.Id)
	}
	return ids
}

// QueryAll 查询所有主题及订阅者ID
func (b *Bus) QueryAll() map[string][]uint64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make(map[string][]uint64, len(b.subscribers))
	for topic, subs := range b.subscribers {
		ids := make([]uint64, 0, len(subs))
		for _, s := range subs {
			ids = append(ids, s.Id)
		}
		result[topic] = ids
	}

	return result
}

// Topics 获取所有主题
func (b *Bus) Topics() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	topics := make([]string, 0, len(b.subscribers))
	for topic := range b.subscribers {
		topics = append(topics, topic)
	}
	return topics
}

// Count 获取指定主题的订阅者数量
func (b *Bus) Count(topic string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.subscribers[topic])
}

// Total 获取所有订阅者总数
func (b *Bus) Total() int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	total := 0
	for _, subs := range b.subscribers {
		total += len(subs)
	}
	return total
}

// Clear 清空所有订阅
func (b *Bus) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subscribers = make(map[string][]*busSubscriber)
}

// ClearTopic 清空指定主题的订阅
func (b *Bus) ClearTopic(topic string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.subscribers, topic)
}
