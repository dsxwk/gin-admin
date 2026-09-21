package eventbus

import (
	"context"
	"log"
	"runtime"
	"sync"
)

const defaultQueueSize = 1024

type config struct {
	workers      int
	queueSize    int
	panicHandler PanicHandler
}

// Handler 泛型事件处理函数
type Handler[T any] func(context.Context, T)

// PanicHandler 事件处理panic回调
type PanicHandler func(topic string, event any, recovered any)

// Option 事件总线配置
type Option func(*config)

// WithWorkers 设置异步处理worker数量
func WithWorkers(count int) Option {
	return func(cfg *config) {
		if count > 0 {
			cfg.workers = count
		}
	}
}

// WithQueueSize 设置异步处理队列大小
func WithQueueSize(size int) Option {
	return func(cfg *config) {
		if size > 0 {
			cfg.queueSize = size
		}
	}
}

// WithPanicHandler 设置事件处理panic回调
func WithPanicHandler(handler PanicHandler) Option {
	return func(cfg *config) {
		if handler != nil {
			cfg.panicHandler = handler
		}
	}
}

type subscriber struct {
	id     uint64
	topic  string
	handle func(context.Context, any)
	async  bool
}

type job struct {
	ctx        context.Context
	subscriber *subscriber
	event      any
}

// Subscription 事件订阅
type Subscription struct {
	bus   *Bus
	topic string
	id    uint64
}

// Unsubscribe 取消事件订阅
func (s *Subscription) Unsubscribe() bool {
	if s == nil || s.bus == nil {
		return false
	}

	return s.bus.unsubscribe(s.topic, s.id)
}

// Bus 底层发布订阅总线
type Bus struct {
	mu           sync.RWMutex
	subscribers  map[string][]*subscriber
	nextID       uint64
	panicHandler PanicHandler

	queue     chan job
	done      chan struct{}
	stopped   chan struct{}
	closeOnce sync.Once
	wg        sync.WaitGroup
}

var (
	defaultBus  *Bus
	defaultOnce sync.Once
)

// Default 获取默认事件总线
func Default() *Bus {
	defaultOnce.Do(func() {
		defaultBus = NewBus()
	})

	return defaultBus
}

// NewBus 创建事件总线
func NewBus(options ...Option) *Bus {
	cfg := config{
		workers:   max(runtime.GOMAXPROCS(0), 1),
		queueSize: defaultQueueSize,
		panicHandler: func(topic string, event any, recovered any) {
			log.Printf("eventbus handler panic: topic=%s event=%T panic=%v", topic, event, recovered)
		},
	}

	for _, option := range options {
		if option != nil {
			option(&cfg)
		}
	}

	bus := &Bus{
		subscribers:  make(map[string][]*subscriber),
		panicHandler: cfg.panicHandler,
		queue:        make(chan job, cfg.queueSize),
		done:         make(chan struct{}),
		stopped:      make(chan struct{}),
	}

	bus.wg.Add(cfg.workers)
	for index := 0; index < cfg.workers; index++ {
		go bus.worker()
	}

	return bus
}

// Subscribe 订阅同步事件
func (b *Bus) Subscribe[T any](topic string, handler Handler[T]) *Subscription {
	return b.subscribe(topic, handler, false)
}

// SubscribeAsync 订阅异步事件
func (b *Bus) SubscribeAsync[T any](topic string, handler Handler[T]) *Subscription {
	return b.subscribe(topic, handler, true)
}

// subscribe 创建指定模式的订阅
func (b *Bus) subscribe[T any](topic string, handler Handler[T], async bool) *Subscription {
	if b == nil || topic == "" || handler == nil || b.isClosed() {
		return nil
	}

	item := &subscriber{
		topic: topic,
		async: async,
		handle: func(ctx context.Context, data any) {
			event, ok := data.(T)
			if !ok {
				return
			}

			handler(ctx, event)
		},
	}

	b.mu.Lock()
	b.nextID++
	item.id = b.nextID
	b.subscribers[topic] = append(b.subscribers[topic], item)
	b.mu.Unlock()

	return &Subscription{
		bus:   b,
		topic: topic,
		id:    item.id,
	}
}

// unsubscribe 取消指定订阅
func (b *Bus) unsubscribe(topic string, id uint64) bool {
	if b == nil {
		return false
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	subscribers := b.subscribers[topic]
	for index, item := range subscribers {
		if item.id != id {
			continue
		}

		copy(subscribers[index:], subscribers[index+1:])
		subscribers[len(subscribers)-1] = nil
		b.subscribers[topic] = subscribers[:len(subscribers)-1]
		if len(b.subscribers[topic]) == 0 {
			delete(b.subscribers, topic)
		}

		return true
	}

	return false
}

// Publish 发布事件
func (b *Bus) Publish(ctx context.Context, topic string, event any) {
	if b == nil || topic == "" || b.isClosed() {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}

	for _, item := range b.subscribersFor(topic) {
		if ctx.Err() != nil {
			return
		}

		if item.async {
			if !b.enqueue(ctx, item, event) {
				return
			}
			continue
		}

		b.handle(ctx, item, event)
	}
}

// enqueue 添加异步任务
func (b *Bus) enqueue(ctx context.Context, item *subscriber, event any) bool {
	select {
	case <-b.done:
		return false
	default:
	}

	select {
	case b.queue <- job{
		ctx:        ctx,
		subscriber: item,
		event:      event,
	}:
		return true
	case <-ctx.Done():
		return false
	case <-b.done:
		return false
	}
}

// subscribersFor 获取订阅者快照
func (b *Bus) subscribersFor(topic string) []*subscriber {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return append([]*subscriber(nil), b.subscribers[topic]...)
}

// worker 消费异步任务
func (b *Bus) worker() {
	defer b.wg.Done()

	for {
		select {
		case item := <-b.queue:
			b.handle(item.ctx, item.subscriber, item.event)
		case <-b.done:
			b.drain()
			return
		}
	}
}

// drain 消费队列中的剩余任务
func (b *Bus) drain() {
	for {
		select {
		case item := <-b.queue:
			b.handle(item.ctx, item.subscriber, item.event)
		default:
			return
		}
	}
}

// handle 执行订阅者
func (b *Bus) handle(ctx context.Context, item *subscriber, event any) {
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx.Err() != nil {
		return
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			b.reportPanic(item.topic, event, recovered)
		}
	}()

	item.handle(ctx, event)
}

// reportPanic 上报事件处理panic
func (b *Bus) reportPanic(topic string, event any, recovered any) {
	if b == nil || b.panicHandler == nil {
		return
	}

	defer func() {
		_ = recover()
	}()

	b.panicHandler(topic, event, recovered)
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

// Close 关闭事件总线并等待异步任务完成
func (b *Bus) Close(ctx context.Context) error {
	if b == nil || b.done == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	b.closeOnce.Do(func() {
		close(b.done)

		go func() {
			b.wg.Wait()
			close(b.stopped)
		}()
	})

	select {
	case <-b.stopped:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// isClosed 判断事件总线是否关闭
func (b *Bus) isClosed() bool {
	if b == nil || b.done == nil {
		return true
	}

	select {
	case <-b.done:
		return true
	default:
		return false
	}
}
