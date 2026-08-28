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
	Grpc          []map[string]any `json:"Grpc"`
	ListenerEvent []map[string]any `json:"ListenerEvent"`
	Job           []map[string]any `json:"Job"`
}

// TraceStore 追踪存储
type TraceStore struct {
	mu    sync.RWMutex
	store map[string]*TraceData
}

var Store = &TraceStore{
	store: make(map[string]*TraceData),
}

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
)

// Get 获取追踪数据,不存在时创建并存储
func (ts *TraceStore) Get(traceId string) *TraceData {
	ts.mu.Lock()
	defer ts.mu.Unlock()

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
	}
}

// addTraceField 通用方法:获取或创建TraceData,并对指定字段追加数据
func addTraceField(traceId string, data map[string]any, fieldFn func(d *TraceData) *[]map[string]any) {
	if traceId == "" {
		return
	}

	Store.mu.Lock()
	d, ok := Store.store[traceId]
	if !ok {
		d = newTraceData()
		Store.store[traceId] = d
	}
	Store.mu.Unlock()

	d.mu.Lock()
	*fieldFn(d) = append(*fieldFn(d), data)
	d.mu.Unlock()
}

var traceFieldMap = map[TraceField]func(d *TraceData) *[]map[string]any{
	FieldSql:      func(d *TraceData) *[]map[string]any { return &d.Sql },
	FieldCache:    func(d *TraceData) *[]map[string]any { return &d.Cache },
	FieldHttp:     func(d *TraceData) *[]map[string]any { return &d.Http },
	FieldMq:       func(d *TraceData) *[]map[string]any { return &d.Mq },
	FieldGrpc:     func(d *TraceData) *[]map[string]any { return &d.Grpc },
	FieldListener: func(d *TraceData) *[]map[string]any { return &d.ListenerEvent },
	FieldJob:      func(d *TraceData) *[]map[string]any { return &d.Job },
}

// Add 记录调试信息
func Add(traceId string, field TraceField, data map[string]any) {
	fieldFn, ok := traceFieldMap[field]
	if !ok {
		return
	}
	addTraceField(traceId, data, fieldFn)
}
