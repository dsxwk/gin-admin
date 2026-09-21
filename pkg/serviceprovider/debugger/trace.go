package debugger

import (
	"encoding/json"
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

// newTraceData 创建追踪数据
func newTraceData() *TraceData {
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

// Append 追加追踪事件
func (d *TraceData) Append(event any) {
	d.mu.Lock()
	defer d.mu.Unlock()

	switch value := event.(type) {
	case SQLEvent:
		d.SQL = append(d.SQL, value)
	case CacheEvent:
		d.Cache = append(d.Cache, value)
	case HTTPEvent:
		d.HTTP = append(d.HTTP, value)
	case MQEvent:
		d.MQ = append(d.MQ, value)
	case GRPCEvent:
		d.GRPC = append(d.GRPC, value)
	case ListenerEvent:
		d.Listener = append(d.Listener, value)
	case JobEvent:
		d.Job = append(d.Job, value)
	case ESEvent:
		d.ES = append(d.ES, value)
	}
}

// MarshalJSON 序列化追踪数据
func (d *TraceData) MarshalJSON() ([]byte, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	type traceData TraceData
	return json.Marshal((*traceData)(d))
}

// Trace 追踪存储
type Trace struct {
	mu   sync.RWMutex
	data map[string]*TraceData
}

// Store 全局追踪存储
var Store = NewTrace()

// NewTrace 创建追踪存储
func NewTrace() *Trace {
	return &Trace{}
}

// Get 获取追踪数据
func (s *Trace) Get(traceID string) (*TraceData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	trace, ok := s.data[traceID]
	return trace, ok
}

// GetOrCreate 获取或创建追踪数据
func (s *Trace) GetOrCreate(traceID string) *TraceData {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.data == nil {
		s.data = make(map[string]*TraceData)
	}

	trace, ok := s.data[traceID]
	if !ok {
		trace = newTraceData()
		s.data[traceID] = trace
	}

	return trace
}

// Record 记录追踪事件
func (s *Trace) Record(traceID string, event any) {
	if traceID == "" {
		return
	}

	s.GetOrCreate(traceID).Append(event)
}

// Delete 删除追踪数据
func (s *Trace) Delete(traceID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, traceID)
}

// CleanExpired 清理过期追踪数据
func (s *Trace) CleanExpired(ttl time.Duration) int {
	if ttl <= 0 {
		return 0
	}

	now := time.Now()
	count := 0

	s.mu.Lock()
	defer s.mu.Unlock()

	for traceID, trace := range s.data {
		if now.Sub(trace.createdAt) < ttl {
			continue
		}

		delete(s.data, traceID)
		count++
	}

	return count
}

// StartCleanup 启动定时清理
func (s *Trace) StartCleanup(interval, ttl time.Duration, stop <-chan struct{}) {
	if interval <= 0 || ttl <= 0 {
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
