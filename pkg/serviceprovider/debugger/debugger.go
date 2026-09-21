package debugger

import (
	"context"
	"gin/common/ctxkey"
	"gin/pkg/serviceprovider/eventbus"
	"sync"
	"time"
)

const (
	traceCleanupInterval = time.Minute
	traceExpire          = 30 * time.Minute
)

var debugTopics = [...]string{
	TopicSQL,
	TopicCache,
	TopicHTTP,
	TopicMQ,
	TopicGRPC,
	TopicListener,
	TopicJob,
	TopicES,
}

// Debugger 调试器入口
type Debugger struct {
	mu            sync.RWMutex
	bus           *eventbus.Bus
	store         *Trace
	subscriptions []*eventbus.Subscription
	cleanupStop   chan struct{}
}

// New 创建调试器
func New(bus *eventbus.Bus) *Debugger {
	return WithStore(bus, Store)
}

// WithStore 使用指定存储创建调试器
func WithStore(bus *eventbus.Bus, store *Trace) *Debugger {
	if bus == nil {
		bus = eventbus.Default()
	}
	if store == nil {
		store = NewTrace()
	}

	return &Debugger{
		bus:   bus,
		store: store,
	}
}

// Start 启动调试器
func (d *Debugger) Start() {
	if d == nil {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cleanupStop != nil {
		return
	}

	d.subscribe()
	d.cleanupStop = make(chan struct{})
	d.store.StartCleanup(traceCleanupInterval, traceExpire, d.cleanupStop)
}

// Stop 停止调试器
func (d *Debugger) Stop() {
	if d == nil {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.cleanupStop == nil {
		return
	}

	d.unsubscribe()
	close(d.cleanupStop)
	d.cleanupStop = nil
}

// Store 获取调试器使用的追踪存储
func (d *Debugger) Store() *Trace {
	if d == nil {
		return nil
	}

	return d.store
}

// IsRunning 判断调试器是否运行中
func (d *Debugger) IsRunning() bool {
	if d == nil {
		return false
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.cleanupStop != nil
}

// subscribe 订阅全部调试事件
func (d *Debugger) subscribe() {
	for _, topic := range debugTopics {
		d.subscriptions = append(d.subscriptions, d.bus.Subscribe[any](topic, d.handleEvent))
	}
}

// unsubscribe 取消全部调试事件
func (d *Debugger) unsubscribe() {
	for _, subscription := range d.subscriptions {
		subscription.Unsubscribe()
	}

	d.subscriptions = nil
}

// handleEvent 处理调试事件
func (d *Debugger) handleEvent(_ context.Context, event any) {
	switch value := event.(type) {
	case SQLEvent:
		d.store.Record(value.TraceID, value)
	case CacheEvent:
		d.store.Record(value.TraceID, value)
	case HTTPEvent:
		d.store.Record(value.TraceID, value)
	case MQEvent:
		d.store.Record(value.TraceID, value)
	case GRPCEvent:
		d.store.Record(value.TraceID, value)
	case JobEvent:
		d.store.Record(value.TraceID, value)
	case ESEvent:
		d.store.Record(value.TraceID, value)
	case eventbus.PublishedEvent:
		traceID := ctxkey.TraceID(value.Context)

		d.store.Record(traceID, ListenerEvent{
			TraceID:     traceID,
			Name:        value.Name,
			Description: value.Description,
			Data:        value.Data,
		})
	}
}
