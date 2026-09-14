package debugger

import (
	"gin/pkg/serviceprovider/eventbus"
	"sync"
	"time"
)

const (
	traceCleanupInterval = time.Minute
	traceExpire          = 30 * time.Minute
)

// Debugger 调试器入口
type Debugger struct {
	mu          sync.RWMutex
	bus         *eventbus.Bus
	store       *TraceStore
	collector   *Collector
	cleanupStop chan struct{}
}

// New 创建调试器
func New(bus *eventbus.Bus) *Debugger {
	return NewWithStore(bus, Store)
}

// NewWithStore 使用指定存储创建调试器
func NewWithStore(bus *eventbus.Bus, store *TraceStore) *Debugger {
	if bus == nil {
		bus = eventbus.NewBus()
	}
	if store == nil {
		store = NewTraceStore()
	}

	return &Debugger{
		bus:       bus,
		store:     store,
		collector: NewCollector(bus, store),
	}
}

// Start 启动调试器
func (d *Debugger) Start() {
	if d == nil {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.collector.IsRunning() {
		return
	}

	d.collector.Register()
	if d.cleanupStop == nil {
		d.cleanupStop = make(chan struct{})
		d.store.StartCleanup(traceCleanupInterval, traceExpire, d.cleanupStop)
	}
}

// Stop 停止调试器
func (d *Debugger) Stop() {
	if d == nil {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	d.collector.Unregister()
	if d.cleanupStop != nil {
		close(d.cleanupStop)
		d.cleanupStop = nil
	}
}

// Bus 获取调试器使用的总线
func (d *Debugger) Bus() *eventbus.Bus {
	if d == nil {
		return nil
	}

	return d.bus
}

// Store 获取调试器使用的追踪存储
func (d *Debugger) Store() *TraceStore {
	if d == nil {
		return nil
	}

	return d.store
}

// SubIds 获取全部订阅ID
func (d *Debugger) SubIds() map[string]uint64 {
	if d == nil {
		return nil
	}

	return d.collector.SubIDs()
}

// IsRunning 判断调试器是否运行中
func (d *Debugger) IsRunning() bool {
	if d == nil {
		return false
	}

	return d.collector.IsRunning()
}

// GetSubId 获取指定主题的订阅ID
func (d *Debugger) GetSubId(topic string) (uint64, bool) {
	if d == nil {
		return 0, false
	}

	return d.collector.GetSubID(topic)
}
