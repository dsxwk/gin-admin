package message

import (
	"sync"
	"sync/atomic"
)

var (
	event     *Event
	eventOnce sync.Once
)

func NewEvent() *Event {
	eventOnce.Do(func() {
		event = &Event{
			subscribers: make(map[string][]*subscriber),
		}
	})

	return event
}

type SubscriberFunc func(event any)

type subscriber struct {
	Id     uint64
	Async  bool
	Handle SubscriberFunc
}

type Event struct {
	subscribers map[string][]*subscriber
	mu          sync.RWMutex
	idCounter   uint64
}

// SubscribeAsync 订阅(异步)
func (b *Event) SubscribeAsync(topic string, fn SubscriberFunc) uint64 {
	return b.addSubscriber(topic, fn, true)
}

// Subscribe 订阅(同步)
func (b *Event) Subscribe(topic string, fn SubscriberFunc) uint64 {
	return b.addSubscriber(topic, fn, false)
}

// 通用订阅
func (b *Event) addSubscriber(topic string, fn SubscriberFunc, async bool) uint64 {
	id := atomic.AddUint64(&b.idCounter, 1)

	sub := &subscriber{
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
func (b *Event) Unsubscribe(topic string, id uint64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	subs, ok := b.subscribers[topic]
	if !ok {
		return false
	}

	newSubs := subs[:0] // 原地过滤
	removed := false

	for _, s := range subs {
		if s.Id == id {
			removed = true
			continue
		}
		newSubs = append(newSubs, s)
	}

	if removed {
		b.subscribers[topic] = newSubs
	}

	return removed
}

// Publish 发布事件
func (b *Event) Publish(topic string, event any) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, sub := range b.subscribers[topic] {
		if sub.Async {
			// 异步执行
			go sub.Handle(event)
		} else {
			// 同步执行
			sub.Handle(event)
		}
	}
}

// SubscribeIds 查询该topic所有订阅者id
func (b *Event) SubscribeIds(topic string) []uint64 {
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

// QueryAll 查询所有topic+订阅者Id
func (b *Event) QueryAll() map[string][]uint64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make(map[string][]uint64)

	for topic, subs := range b.subscribers {
		ids := make([]uint64, 0, len(subs))
		for _, s := range subs {
			ids = append(ids, s.Id)
		}
		result[topic] = ids
	}

	return result
}
