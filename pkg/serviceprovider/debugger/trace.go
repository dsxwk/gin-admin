package debugger

import (
	"encoding/json"
	"maps"
	"sync"
	"time"
)

// TraceData 单个追踪数据
type TraceData struct {
	mu sync.RWMutex

	SQL      []SQLEvent      `json:"sql"`
	Cache    []CacheEvent    `json:"cache"`
	HTTP     []HTTPEvent     `json:"http"`
	MQ       []MQEvent       `json:"mq"`
	GRPC     []GRPCEvent     `json:"grpc"`
	Listener []ListenerEvent `json:"listener"`
	Job      []JobEvent      `json:"job"`
	ES       []ESEvent       `json:"es"`

	createdAt time.Time
}

// NewTraceData 创建追踪数据
func NewTraceData() *TraceData {
	return &TraceData{
		SQL:       make([]SQLEvent, 0),
		Cache:     make([]CacheEvent, 0),
		HTTP:      make([]HTTPEvent, 0),
		MQ:        make([]MQEvent, 0),
		GRPC:      make([]GRPCEvent, 0),
		Listener:  make([]ListenerEvent, 0),
		Job:       make([]JobEvent, 0),
		ES:        make([]ESEvent, 0),
		createdAt: time.Now(),
	}
}

// CreatedAt 获取创建时间
func (d *TraceData) CreatedAt() time.Time {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return d.createdAt
}

// AddSQL 添加SQL事件
func (d *TraceData) AddSQL(event SQLEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.SQL = append(d.SQL, event)
}

// AddCache 添加缓存事件
func (d *TraceData) AddCache(event CacheEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.Cache = append(d.Cache, event)
}

// AddHTTP 添加HTTP事件
func (d *TraceData) AddHTTP(event HTTPEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.HTTP = append(d.HTTP, event)
}

// AddMQ 添加消息队列事件
func (d *TraceData) AddMQ(event MQEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.MQ = append(d.MQ, event)
}

// AddGRPC 添加gRPC事件
func (d *TraceData) AddGRPC(event GRPCEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.GRPC = append(d.GRPC, event)
}

// AddListener 添加业务监听事件
func (d *TraceData) AddListener(event ListenerEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.Listener = append(d.Listener, event)
}

// AddJob 添加Job事件
func (d *TraceData) AddJob(event JobEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.Job = append(d.Job, event)
}

// AddES 添加ES事件
func (d *TraceData) AddES(event ESEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.ES = append(d.ES, event)
}

// MarshalJSON 序列化追踪数据
func (d *TraceData) MarshalJSON() ([]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	type traceData TraceData
	return json.Marshal((*traceData)(d))
}

// Summary 获取统计摘要
func (d *TraceData) Summary() map[string]int {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return map[string]int{
		"sql":      len(d.SQL),
		"cache":    len(d.Cache),
		"http":     len(d.HTTP),
		"mq":       len(d.MQ),
		"grpc":     len(d.GRPC),
		"listener": len(d.Listener),
		"job":      len(d.Job),
		"es":       len(d.ES),
	}
}

// TraceStore 追踪存储
type TraceStore struct {
	mu   sync.RWMutex
	data map[string]*TraceData
}

// Store 全局追踪存储
var Store = NewTraceStore()

// NewTraceStore 创建追踪存储
func NewTraceStore() *TraceStore {
	return &TraceStore{
		data: make(map[string]*TraceData),
	}
}

// Get 获取追踪数据
func (s *TraceStore) Get(traceID string) (*TraceData, bool) {
	if s == nil {
		return nil, false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	trace, ok := s.data[traceID]
	return trace, ok
}

// GetOrCreate 获取或创建追踪数据
func (s *TraceStore) GetOrCreate(traceID string) *TraceData {
	if s == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[string]*TraceData)
	}

	trace, ok := s.data[traceID]
	if ok {
		return trace
	}

	trace = NewTraceData()
	s.data[traceID] = trace
	return trace
}

// Set 设置追踪数据
func (s *TraceStore) Set(traceID string, trace *TraceData) {
	if s == nil || trace == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[string]*TraceData)
	}
	s.data[traceID] = trace
}

// Delete 删除追踪数据
func (s *TraceStore) Delete(traceID string) {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, traceID)
}

// Has 判断追踪数据是否存在
func (s *TraceStore) Has(traceID string) bool {
	if s == nil {
		return false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.data[traceID]
	return ok
}

// Count 获取追踪数据数量
func (s *TraceStore) Count() int {
	if s == nil {
		return 0
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.data)
}

// List 获取全部追踪数据
func (s *TraceStore) List() map[string]*TraceData {
	if s == nil {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*TraceData, len(s.data))
	maps.Copy(result, s.data)

	return result
}

// Clear 清空全部追踪数据
func (s *TraceStore) Clear() {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = make(map[string]*TraceData)
}

// CleanExpired 清理过期追踪数据
func (s *TraceStore) CleanExpired(ttl time.Duration) int {
	if s == nil || ttl <= 0 {
		return 0
	}

	now := time.Now()
	count := 0

	s.mu.Lock()
	defer s.mu.Unlock()

	for traceID, trace := range s.data {
		if now.Sub(trace.CreatedAt()) < ttl {
			continue
		}

		delete(s.data, traceID)
		count++
	}

	return count
}

// StartCleanup 启动定时清理
func (s *TraceStore) StartCleanup(interval, ttl time.Duration, stop <-chan struct{}) {
	if s == nil || interval <= 0 || ttl <= 0 {
		return
	}

	ticker := time.NewTicker(interval)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.CleanExpired(ttl)
			case <-stop:
				return
			}
		}
	}()
}
