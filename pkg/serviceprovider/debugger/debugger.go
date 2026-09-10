package debugger

import (
	"gin/pkg/serviceprovider/message"
	"sync"
)

type Debugger struct {
	Bus    *message.Event
	subIds map[string]uint64
	mu     sync.RWMutex // 读写锁
}

func NewDebugger(bus *message.Event) *Debugger {
	return &Debugger{
		Bus:    bus,
		subIds: make(map[string]uint64),
	}
}

func (d *Debugger) Start() {
	if d.Bus == nil {
		return
	}
	id1 := d.Bus.Subscribe(TopicSql, func(e SqlEvent) {
		Add(e.TraceId, FieldSql, map[string]any{
			"sql":  e.Sql,
			"rows": e.Rows,
			"ms":   e.Ms,
		})
	})
	id2 := d.Bus.Subscribe(TopicCache, func(e CacheEvent) {
		Add(e.TraceId, FieldCache, map[string]any{
			"driver": e.Driver,
			"name":   e.Name,
			"cmd":    e.Cmd,
			"args":   e.Args,
			"ms":     e.Ms,
		})
	})
	id3 := d.Bus.Subscribe(TopicHttp, func(e HttpEvent) {
		Add(e.TraceId, FieldHttp, map[string]any{
			"url":      e.Url,
			"method":   e.Method,
			"header":   e.Header,
			"body":     e.Body,
			"status":   e.Status,
			"response": e.Response,
			"ms":       e.Ms,
		})
	})
	id4 := d.Bus.Subscribe(TopicMq, func(e MqEvent) {
		Add(e.TraceId, FieldMq, map[string]any{
			"driver":  e.Driver,
			"topic":   e.Topic,
			"message": e.Message,
			"key":     e.Key,
			"group":   e.Group,
			"ms":      e.Ms,
			"extra":   e.Extra,
		})
	})
	id5 := d.Bus.Subscribe(TopicGrpc, func(e GrpcEvent) {
		Add(e.TraceId, FieldGrpc, map[string]any{
			"method":   e.Method,
			"request":  e.Request,
			"response": e.Response,
			"code":     e.Code,
			"ms":       e.Ms,
		})
	})
	id6 := d.Bus.Subscribe(TopicListener, func(e ListenerEvent) {
		Add(e.TraceId, FieldListener, map[string]any{
			"name":  e.Name,
			"topic": e.Description,
			"data":  e.Data,
		})
	})
	id7 := d.Bus.Subscribe(TopicJob, func(e JobEvent) {
		Add(e.TraceId, FieldJob, map[string]any{
			"name":       e.Name,
			"connection": e.Connection,
			"payload":    e.Payload,
			"ms":         e.Ms,
		})
	})
	id8 := d.Bus.Subscribe(TopicEs, func(e EsEvent) {
		Add(e.TraceId, FieldEs, map[string]any{
			"action":   e.Action,
			"index":    e.Index,
			"request":  e.Request,
			"response": e.Response,
			"code":     e.Code,
			"ms":       e.Ms,
		})
	})

	d.mu.Lock()
	defer d.mu.Unlock()

	d.subIds[TopicSql] = id1
	d.subIds[TopicCache] = id2
	d.subIds[TopicHttp] = id3
	d.subIds[TopicMq] = id4
	d.subIds[TopicGrpc] = id5
	d.subIds[TopicListener] = id6
	d.subIds[TopicJob] = id7
	d.subIds[TopicEs] = id8
}

func (d *Debugger) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()

	for topic, id := range d.subIds {
		d.Bus.Unsubscribe(topic, id)
		// 清空订阅ID
		delete(d.subIds, topic)
	}
}

// SubIds 获取所有订阅ID(用于调试和检查)
func (d *Debugger) SubIds() map[string]uint64 {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// 返回副本避免外部修改
	result := make(map[string]uint64, len(d.subIds))
	for k, v := range d.subIds {
		result[k] = v
	}
	return result
}

// IsRunning 检查调试器是否运行中
func (d *Debugger) IsRunning() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return len(d.subIds) > 0
}

// GetSubId 获取指定主题的订阅ID
func (d *Debugger) GetSubId(topic string) (uint64, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	id, ok := d.subIds[topic]
	return id, ok
}
