package debugger

import (
	"sync"
)

// TraceData 单个追踪数据
type TraceData struct {
	mu            sync.RWMutex
	Sql           []map[string]any `json:"Sql"`
	Cache         []map[string]any `json:"Cache"`
	Http          []map[string]any `json:"Http"`
	Mq            []map[string]any `json:"Mq"`
	ListenerEvent []map[string]any `json:"ListenerEvent"`
}

// TraceStore 追踪存储
type TraceStore struct {
	mu    sync.RWMutex
	store map[string]*TraceData
}

var Store = &TraceStore{
	store: make(map[string]*TraceData),
}

// Get 获取追踪数据,不存在时创建并存储
func (ts *TraceStore) Get(traceId string) *TraceData {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if data, ok := ts.store[traceId]; ok {
		return data
	}

	data := &TraceData{
		Sql:           make([]map[string]any, 0),
		Cache:         make([]map[string]any, 0),
		Http:          make([]map[string]any, 0),
		Mq:            make([]map[string]any, 0),
		ListenerEvent: make([]map[string]any, 0),
	}
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

// addTraceField 通用方法:获取或创建TraceData,并对指定字段追加数据
func addTraceField(traceId string, data map[string]any, fieldFn func(d *TraceData) *[]map[string]any) {
	if traceId == "" {
		return
	}

	Store.mu.Lock()
	d, ok := Store.store[traceId]
	if !ok {
		d = &TraceData{
			Sql:           make([]map[string]any, 0),
			Cache:         make([]map[string]any, 0),
			Http:          make([]map[string]any, 0),
			Mq:            make([]map[string]any, 0),
			ListenerEvent: make([]map[string]any, 0),
		}
		Store.store[traceId] = d
	}
	Store.mu.Unlock()

	d.mu.Lock()
	*fieldFn(d) = append(*fieldFn(d), data)
	d.mu.Unlock()
}

// AddSql 记录sql调试信息
func AddSql(traceId string, data map[string]any) {
	addTraceField(traceId, data, func(d *TraceData) *[]map[string]any { return &d.Sql })
}

// AddCache 记录缓存调试信息
func AddCache(traceId string, data map[string]any) {
	addTraceField(traceId, data, func(d *TraceData) *[]map[string]any { return &d.Cache })
}

// AddHttp 记录http调试信息
func AddHttp(traceId string, data map[string]any) {
	addTraceField(traceId, data, func(d *TraceData) *[]map[string]any { return &d.Http })
}

// AddMq 记录mq调试信息
func AddMq(traceId string, data map[string]any) {
	addTraceField(traceId, data, func(d *TraceData) *[]map[string]any { return &d.Mq })
}

// AddListener 记录监听调试信息
func AddListener(traceId string, data map[string]any) {
	addTraceField(traceId, data, func(d *TraceData) *[]map[string]any { return &d.ListenerEvent })
}
