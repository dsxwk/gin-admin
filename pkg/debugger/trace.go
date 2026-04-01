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
	store map[string]*TraceData // 关键：存储指针
}

var Store = &TraceStore{
	store: make(map[string]*TraceData),
}

// Get 获取追踪数据
func (ts *TraceStore) Get(traceId string) *TraceData {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	if data, ok := ts.store[traceId]; ok {
		return data
	}

	// 如果不存在,返回一个新的TraceData(但不会自动存入store)
	// 注意：这里返回的是新对象,不会自动存储到map中
	return &TraceData{
		Sql:           make([]map[string]any, 0),
		Cache:         make([]map[string]any, 0),
		Http:          make([]map[string]any, 0),
		Mq:            make([]map[string]any, 0),
		ListenerEvent: make([]map[string]any, 0),
	}
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

// AddSql 记录sql调试信息
func AddSql(traceId string, data map[string]any) {
	if traceId == "" {
		return
	}

	Store.mu.Lock()
	defer Store.mu.Unlock()

	// 获取或创建TraceData
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

	// 直接操作d(指针)
	d.mu.Lock()
	d.Sql = append(d.Sql, data)
	d.mu.Unlock()
}

// AddCache 记录缓存调试信息
func AddCache(traceId string, data map[string]any) {
	if traceId == "" {
		return
	}

	Store.mu.Lock()
	defer Store.mu.Unlock()

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

	d.mu.Lock()
	d.Cache = append(d.Cache, data)
	d.mu.Unlock()
}

// AddHttp 记录http调试信息
func AddHttp(traceId string, data map[string]any) {
	if traceId == "" {
		return
	}

	Store.mu.Lock()
	defer Store.mu.Unlock()

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

	d.mu.Lock()
	d.Http = append(d.Http, data)
	d.mu.Unlock()
}

// AddMq 记录mq调试信息
func AddMq(traceId string, data map[string]any) {
	if traceId == "" {
		return
	}

	Store.mu.Lock()
	defer Store.mu.Unlock()

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

	d.mu.Lock()
	d.Mq = append(d.Mq, data)
	d.mu.Unlock()
}

// AddListener 记录监听调试信息
func AddListener(traceId string, data map[string]any) {
	if traceId == "" {
		return
	}

	Store.mu.Lock()
	defer Store.mu.Unlock()

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

	d.mu.Lock()
	d.ListenerEvent = append(d.ListenerEvent, data)
	d.mu.Unlock()
}
