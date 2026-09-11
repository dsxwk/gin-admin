package debugger

import (
	"sync"
	"time"
)

// TraceField 追踪字段类型
type TraceField string

const (
	FieldSql      TraceField = "Sql"
	FieldCache    TraceField = "Cache"
	FieldHttp     TraceField = "Http"
	FieldMq       TraceField = "Mq"
	FieldGrpc     TraceField = "Grpc"
	FieldListener TraceField = "Listener"
	FieldJob      TraceField = "Job"
	FieldEs       TraceField = "Es"
)

// TraceData 单个追踪数据
type TraceData struct {
	mu            sync.RWMutex
	Sql           []map[string]any `json:"Sql"`
	Cache         []map[string]any `json:"Cache"`
	Http          []map[string]any `json:"Http"`
	Mq            []map[string]any `json:"Mq"`
	Grpc          []map[string]any `json:"Grpc"`
	ListenerEvent []map[string]any `json:"ListenerEvent"`
	Job           []map[string]any `json:"Job"`
	Es            []map[string]any `json:"Es"`
	createdAt     time.Time        // 创建时间(用于过期清理)
}

// newTraceData 创建空的追踪数据
func newTraceData() *TraceData {
	return &TraceData{
		Sql:           make([]map[string]any, 0),
		Cache:         make([]map[string]any, 0),
		Http:          make([]map[string]any, 0),
		Mq:            make([]map[string]any, 0),
		Grpc:          make([]map[string]any, 0),
		ListenerEvent: make([]map[string]any, 0),
		Job:           make([]map[string]any, 0),
		Es:            make([]map[string]any, 0),
	}
}

// traceFieldMap 字段映射
var traceFieldMap = map[TraceField]func(d *TraceData) *[]map[string]any{
	FieldSql:      func(d *TraceData) *[]map[string]any { return &d.Sql },
	FieldCache:    func(d *TraceData) *[]map[string]any { return &d.Cache },
	FieldHttp:     func(d *TraceData) *[]map[string]any { return &d.Http },
	FieldMq:       func(d *TraceData) *[]map[string]any { return &d.Mq },
	FieldGrpc:     func(d *TraceData) *[]map[string]any { return &d.Grpc },
	FieldListener: func(d *TraceData) *[]map[string]any { return &d.ListenerEvent },
	FieldJob:      func(d *TraceData) *[]map[string]any { return &d.Job },
	FieldEs:       func(d *TraceData) *[]map[string]any { return &d.Es },
}

// TraceStore 追踪存储
type TraceStore struct {
	mu    sync.RWMutex
	store map[string]*TraceData
}

// Store 全局追踪存储
var Store = &TraceStore{
	store: make(map[string]*TraceData),
}

// Get 获取追踪数据（不存在时创建）
func (ts *TraceStore) Get(traceId string) *TraceData {
	ts.mu.RLock()
	if data, ok := ts.store[traceId]; ok {
		ts.mu.RUnlock()
		return data
	}
	ts.mu.RUnlock()

	ts.mu.Lock()
	defer ts.mu.Unlock()

	// 双重检查
	if data, ok := ts.store[traceId]; ok {
		return data
	}

	data := newTraceData()
	ts.store[traceId] = data
	return data
}

// Set 设置追踪数据
func (ts *TraceStore) Set(traceId string, data *TraceData) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.store[traceId] = data
}

// Delete 删除追踪数据
func (ts *TraceStore) Delete(traceId string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	delete(ts.store, traceId)
}

// Has 检查是否存在
func (ts *TraceStore) Has(traceId string) bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	_, ok := ts.store[traceId]
	return ok
}

// Count 获取追踪数量
func (ts *TraceStore) Count() int {
	ts.mu.RLock()
	defer ts.mu.RUnlock()
	return len(ts.store)
}

// List 获取所有traceId
func (ts *TraceStore) List() []string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	ids := make([]string, 0, len(ts.store))
	for id := range ts.store {
		ids = append(ids, id)
	}
	return ids
}

// Clear 清空所有追踪数据
func (ts *TraceStore) Clear() {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.store = make(map[string]*TraceData)
}

// CleanExpired 清理过期的追踪数据
func (ts *TraceStore) CleanExpired(expire time.Duration) int {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	now := time.Now()
	cleaned := 0

	for id, data := range ts.store {
		if now.Sub(data.createdAt) > expire {
			delete(ts.store, id)
			cleaned++
		}
	}

	return cleaned
}

// StartCleanup 启动定时清理(返回停止函数)
func (ts *TraceStore) StartCleanup(interval, expire time.Duration) func() {
	ticker := time.NewTicker(interval)
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				ts.CleanExpired(expire)
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		close(stop)
	}
}

// Add 记录调试信息
func Add(traceId string, field TraceField, data map[string]any) {
	if traceId == "" {
		return
	}

	fieldFn, ok := traceFieldMap[field]
	if !ok {
		return
	}

	// 获取或创建TraceData
	Store.mu.Lock()
	d, ok := Store.store[traceId]
	if !ok {
		d = newTraceData()
		Store.store[traceId] = d
	}
	Store.mu.Unlock()

	// 追加数据(只在TraceData锁内操作)
	d.mu.Lock()
	*fieldFn(d) = append(*fieldFn(d), data)
	d.mu.Unlock()
}

// AddBatch 批量记录
func AddBatch(traceId string, items map[TraceField]map[string]any) {
	if traceId == "" || len(items) == 0 {
		return
	}

	// 获取或创建TraceData
	Store.mu.Lock()
	d, ok := Store.store[traceId]
	if !ok {
		d = newTraceData()
		Store.store[traceId] = d
	}
	Store.mu.Unlock()

	// 批量追加
	d.mu.Lock()
	defer d.mu.Unlock()

	for field, data := range items {
		if fieldFn, _ok := traceFieldMap[field]; _ok {
			*fieldFn(d) = append(*fieldFn(d), data)
		}
	}
}

// GetField 获取指定字段的数据
func (d *TraceData) GetField(field TraceField) []map[string]any {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if fieldFn, ok := traceFieldMap[field]; ok {
		src := *fieldFn(d)
		result := make([]map[string]any, len(src))
		copy(result, src)
		return result
	}
	return nil
}

// Summary 获取统计摘要
func (d *TraceData) Summary() map[string]int {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return map[string]int{
		"Sql":           len(d.Sql),
		"Cache":         len(d.Cache),
		"Http":          len(d.Http),
		"Mq":            len(d.Mq),
		"Grpc":          len(d.Grpc),
		"ListenerEvent": len(d.ListenerEvent),
		"Job":           len(d.Job),
		"Es":            len(d.Es),
	}
}
