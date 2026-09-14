package debugger

import (
	"gin/common/ctxkey"
	"gin/pkg/serviceprovider/eventbus"
	"maps"
	"sync"
)

// Collector 调试事件收集器
type Collector struct {
	mu      sync.RWMutex
	bus     *eventbus.Bus
	store   *TraceStore
	subIDs  map[string]uint64
	running bool
}

// NewCollector 创建调试事件收集器
func NewCollector(bus *eventbus.Bus, store *TraceStore) *Collector {
	if bus == nil {
		bus = eventbus.NewBus()
	}
	if store == nil {
		store = NewTraceStore()
	}

	return &Collector{
		bus:    bus,
		store:  store,
		subIDs: make(map[string]uint64),
	}
}

// Register 注册全部调试事件订阅
func (c *Collector) Register() {
	if c == nil || c.bus == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return
	}

	c.subIDs[TopicSQL] = c.bus.Subscribe(TopicSQL, func(event SQLEvent) {
		c.addSQL(event)
	})
	c.subIDs[TopicCache] = c.bus.Subscribe(TopicCache, func(event CacheEvent) {
		c.addCache(event)
	})
	c.subIDs[TopicHTTP] = c.bus.Subscribe(TopicHTTP, func(event HTTPEvent) {
		c.addHTTP(event)
	})
	c.subIDs[TopicMQ] = c.bus.Subscribe(TopicMQ, func(event MQEvent) {
		c.addMQ(event)
	})
	c.subIDs[TopicGRPC] = c.bus.Subscribe(TopicGRPC, func(event GRPCEvent) {
		c.addGRPC(event)
	})
	c.subIDs[TopicListener] = c.bus.Subscribe(TopicListener, func(event eventbus.PublishedEvent) {
		c.addBusinessEvent(event)
	})
	c.subIDs[TopicJob] = c.bus.Subscribe(TopicJob, func(event JobEvent) {
		c.addJob(event)
	})
	c.subIDs[TopicES] = c.bus.Subscribe(TopicES, func(event ESEvent) {
		c.addES(event)
	})

	c.running = true
}

// Unregister 取消全部调试事件订阅
func (c *Collector) Unregister() {
	if c == nil || c.bus == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for topic, id := range c.subIDs {
		c.bus.Unsubscribe(topic, id)
		delete(c.subIDs, topic)
	}
	c.running = false
}

// IsRunning 判断收集器是否运行中
func (c *Collector) IsRunning() bool {
	if c == nil {
		return false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.running
}

// SubIDs 获取全部订阅ID
func (c *Collector) SubIDs() map[string]uint64 {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]uint64, len(c.subIDs))
	maps.Copy(result, c.subIDs)

	return result
}

// GetSubID 获取指定主题的订阅ID
func (c *Collector) GetSubID(topic string) (uint64, bool) {
	if c == nil {
		return 0, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	id, ok := c.subIDs[topic]
	return id, ok
}

func (c *Collector) addSQL(event SQLEvent) {
	if event.TraceID == "" {
		return
	}

	c.store.GetOrCreate(event.TraceID).AddSQL(event)
}

func (c *Collector) addCache(event CacheEvent) {
	if event.TraceID == "" {
		return
	}

	c.store.GetOrCreate(event.TraceID).AddCache(event)
}

func (c *Collector) addHTTP(event HTTPEvent) {
	if event.TraceID == "" {
		return
	}

	c.store.GetOrCreate(event.TraceID).AddHTTP(event)
}

func (c *Collector) addMQ(event MQEvent) {
	if event.TraceID == "" {
		return
	}

	c.store.GetOrCreate(event.TraceID).AddMQ(event)
}

func (c *Collector) addGRPC(event GRPCEvent) {
	if event.TraceID == "" {
		return
	}

	c.store.GetOrCreate(event.TraceID).AddGRPC(event)
}

func (c *Collector) addBusinessEvent(event eventbus.PublishedEvent) {
	traceID := ctxkey.GetTraceId(event.Context)
	if traceID == "" {
		return
	}

	c.store.GetOrCreate(traceID).AddListener(ListenerEvent{
		TraceID:     traceID,
		Name:        event.Name,
		Description: event.Description,
		Data:        event.Data,
	})
}

func (c *Collector) addJob(event JobEvent) {
	if event.TraceID == "" {
		return
	}

	c.store.GetOrCreate(event.TraceID).AddJob(event)
}

func (c *Collector) addES(event ESEvent) {
	if event.TraceID == "" {
		return
	}

	c.store.GetOrCreate(event.TraceID).AddES(event)
}
