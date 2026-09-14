package eventbus

import (
	"context"
	"fmt"
	"sync"
)

// TopicEvent 业务事件主题
const TopicEvent = "eventbus:event"

// EventInfo 业务事件信息
type EventInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Listeners   []string `json:"listeners"`
}

// Registry 业务事件注册表
type Registry struct {
	mu    sync.RWMutex
	bus   *Bus
	infos map[string]*EventInfo
}

// NewRegistry 创建业务事件注册表
func NewRegistry(bus *Bus) *Registry {
	if bus == nil {
		bus = NewBus()
	}

	return &Registry{
		bus:   bus,
		infos: make(map[string]*EventInfo),
	}
}

// Bus 获取注册表使用的底层总线
func (r *Registry) Bus() *Bus {
	return r.bus
}

// Register 注册业务事件监听器
func (r *Registry) Register[T Event](listener Listener[T], event T) {
	if r == nil || r.bus == nil || listener == nil {
		return
	}

	name := event.Name()
	r.bus.SubscribeAsync(name, func(data any) {
		value, ok := data.(T)
		if !ok {
			return
		}

		listener.Handle(value)
	})

	r.mu.Lock()
	defer r.mu.Unlock()

	info, ok := r.infos[name]
	if !ok {
		info = &EventInfo{
			Name:        name,
			Description: event.Description(),
		}
		r.infos[name] = info
	}

	info.Listeners = append(info.Listeners, fmt.Sprintf("%T", listener))
}

// Publish 发布业务事件
func (r *Registry) Publish[T Event](ctx context.Context, event T) {
	if r == nil || r.bus == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}

	r.bus.PublishWithContext(ctx, TopicEvent, PublishedEvent{
		Context:     ctx,
		Name:        event.Name(),
		Description: event.Description(),
		Data:        event,
	})
	r.bus.PublishWithContext(ctx, event.Name(), event)
}

// EventList 获取全部业务事件信息
func (r *Registry) EventList() []EventInfo {
	if r == nil {
		return nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]EventInfo, 0, len(r.infos))
	for _, info := range r.infos {
		result = append(result, cloneEventInfo(info))
	}

	return result
}

// cloneEventInfo 复制事件信息
func cloneEventInfo(info *EventInfo) EventInfo {
	result := *info
	result.Listeners = append([]string(nil), info.Listeners...)
	return result
}
