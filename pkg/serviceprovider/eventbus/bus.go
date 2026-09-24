package eventbus

import (
	"context"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultQueueSize = 1024
	// defaultEnqueueWait 队列满时等待空闲槽位的最长时间
	defaultEnqueueWait = 100 * time.Millisecond
)

type config struct {
	workers      int
	queueSize    int
	enqueueWait  time.Duration
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

// WithEnqueueWait 设置队列满时等待空闲槽位的最长时间,0为立即丢弃
func WithEnqueueWait(wait time.Duration) Option {
	return func(cfg *config) {
		if wait < 0 {
			wait = 0
		}

		cfg.enqueueWait = wait
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

	queue       chan job
	enqueueWait time.Duration
	dropped     atomic.Uint64
	mismatched  atomic.Uint64
	done        chan struct{}
	stopped     chan struct{}
	closeOnce   sync.Once
	wg          sync.WaitGroup
}

// Stats 事件总线运行统计
type Stats struct {
	Queue      int    // 当前队列长度
	Capacity   int    // 队列容量
	Dropped    uint64 // 队列满丢弃的事件数
	Mismatched uint64 // 事件类型不匹配次数
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
		workers:     max(runtime.GOMAXPROCS(0), 1),
		queueSize:   defaultQueueSize,
		enqueueWait: defaultEnqueueWait,
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
		enqueueWait:  cfg.enqueueWait,
		done:         make(chan struct{}),
		stopped:      make(chan struct{}),
	}

	for index := 0; index < cfg.workers; index++ {
		// WaitGroup.Go 内部负责Done,worker内不要再调用wg.Done
		bus.wg.Go(bus.worker)
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
				b.reportMismatch(topic, data)
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
		// 异步订阅与发布方context解耦,只受队列容量限制
		if item.async {
			b.enqueue(ctx, item, event)
			continue
		}

		if ctx.Err() != nil {
			return
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

	// 异步任务剥离取消信号但保留value,request结束后仍能执行
	task := job{
		ctx:        context.WithoutCancel(ctx),
		subscriber: item,
		event:      event,
	}

	select {
	case b.queue <- task:
		return true
	default:
	}

	// 队列已满,短暂等待空闲槽位,超时丢弃避免拖垮发布方
	timer := time.NewTimer(b.enqueueWait)
	defer timer.Stop()

	select {
	case b.queue <- task:
		return true
	case <-ctx.Done():
		return false
	case <-b.done:
		return false
	case <-timer.C:
		b.reportDrop(item.topic, event)
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

// reportDrop 上报队列满丢弃事件
func (b *Bus) reportDrop(topic string, event any) {
	if b == nil {
		return
	}

	b.dropped.Add(1)
	log.Printf("eventbus queue full, event dropped: topic=%s event=%T", topic, event)
}

// reportMismatch 上报事件类型不匹配
func (b *Bus) reportMismatch(topic string, event any) {
	if b == nil {
		return
	}

	b.mismatched.Add(1)
	log.Printf("eventbus event type mismatch: topic=%s event=%T", topic, event)
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

// Stats 获取事件总线运行统计
func (b *Bus) Stats() Stats {
	if b == nil {
		return Stats{}
	}

	return Stats{
		Queue:      len(b.queue),
		Capacity:   cap(b.queue),
		Dropped:    b.dropped.Load(),
		Mismatched: b.mismatched.Load(),
	}
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
