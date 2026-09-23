package eventbus

import (
	"context"
	"fmt"
	"slices"
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

// ListenerRegistration 监听器注册项
type ListenerRegistration func(*Registry)

// Registry 业务事件注册表
type Registry struct {
	mu    sync.RWMutex
	bus   *Bus
	infos map[string]*EventInfo
}

// Bind 绑定监听器与事件
func Bind[T Event](listener Listener[T], event T) ListenerRegistration {
	return func(registry *Registry) {
		registry.Register(listener, event)
	}
}

// Register 批量注册监听器
func Register(registry *Registry, registrations []ListenerRegistration) {
	for _, registration := range registrations {
		if registration == nil {
			continue
		}
		registration(registry)
	}
}

// NewRegistry 创建业务事件注册表
func NewRegistry(bus *Bus) *Registry {
	if bus == nil {
		bus = Default()
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

// Register 注册业务事件监听器,返回是否注册成功
func (r *Registry) Register[T Event](listener Listener[T], event T) bool {
	if r == nil || r.bus == nil || listener == nil {
		return false
	}

	name := event.Name()
	if name == "" {
		return false
	}

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

	listenerName := fmt.Sprintf("%T", listener)
	if slices.Contains(info.Listeners, listenerName) {
		return false
	}

	subscription := r.bus.SubscribeAsync(name, func(_ context.Context, data any) {
		value, ok := data.(T)
		if !ok {
			return
		}

		listener.Handle(value)
	})
	if subscription == nil {
		return false
	}
	info.Listeners = append(info.Listeners, listenerName)

	return true
}

// Publish 发布业务事件
func (r *Registry) Publish[T Event](ctx context.Context, event T) {
	if r == nil || r.bus == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}

	r.bus.Publish(ctx, TopicEvent, PublishedEvent{
		Context:     ctx,
		Name:        event.Name(),
		Description: event.Description(),
		Data:        event,
	})
	r.bus.Publish(ctx, event.Name(), event)
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
