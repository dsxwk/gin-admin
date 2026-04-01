package debugger

import (
	"gin/pkg/message"
)

type Debugger struct {
	Bus    *message.EventBus
	subIds map[string]uint64
}

func NewDebugger(bus *message.EventBus) *Debugger {
	return &Debugger{
		Bus:    bus,
		subIds: make(map[string]uint64),
	}
}

func (d *Debugger) Start() {
	if d.Bus == nil {
		return
	}
	d.subIds[TopicSql] = d.Bus.Subscribe(TopicSql, func(ev any) {
		if e, ok := ev.(SqlEvent); ok {
			AddSql(e.TraceId, map[string]any{
				"sql":  e.Sql,
				"rows": e.Rows,
				"ms":   e.Ms,
			})
		}
	})
	d.subIds[TopicCache] = d.Bus.Subscribe(TopicCache, func(ev any) {
		if e, ok := ev.(CacheEvent); ok {
			AddCache(e.TraceId, map[string]any{
				"driver": e.Driver,
				"name":   e.Name,
				"cmd":    e.Cmd,
				"args":   e.Args,
				"ms":     e.Ms,
			})
		}
	})
	d.subIds[TopicHttp] = d.Bus.Subscribe(TopicHttp, func(ev any) {
		if e, ok := ev.(HttpEvent); ok {
			AddHttp(e.TraceId, map[string]any{
				"url":      e.Url,
				"method":   e.Method,
				"header":   e.Header,
				"body":     e.Body,
				"status":   e.Status,
				"response": e.Response,
				"ms":       e.Ms,
			})
		}
	})
	d.subIds[TopicMq] = d.Bus.Subscribe(TopicMq, func(ev any) {
		if e, ok := ev.(MqEvent); ok {
			AddMq(e.TraceId, map[string]any{
				"driver":  e.Driver,
				"topic":   e.Topic,
				"message": e.Message,
				"key":     e.Key,
				"group":   e.Group,
				"ms":      e.Ms,
				"extra":   e.Extra,
			})
		}
	})
	d.subIds[TopicListener] = d.Bus.Subscribe(TopicListener, func(ev any) {
		if e, ok := ev.(ListenerEvent); ok {
			AddListener(e.TraceId, map[string]any{
				"name":  e.Name,
				"topic": e.Description,
				"data":  e.Data,
			})
		}
	})
}

func (d *Debugger) Stop() {
	for topic, id := range d.subIds {
		d.Bus.Unsubscribe(topic, id)
	}
	// 清空订阅ID
	for k := range d.subIds {
		delete(d.subIds, k)
	}
}

// SubIds 获取所有订阅ID(用于调试和检查)
func (d *Debugger) SubIds() map[string]uint64 {
	return d.subIds
}

// IsRunning 检查调试器是否运行中
func (d *Debugger) IsRunning() bool {
	return len(d.subIds) > 0
}
