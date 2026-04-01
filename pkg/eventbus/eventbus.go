package eventbus

import (
	"context"
	"fmt"
	"gin/common/base"
	"gin/common/ctxkey"
	"gin/pkg/debugger"
	"gin/pkg/message"
	"github.com/fatih/color"
	"sync"
)

type EventInfo struct {
	Name        string
	Description string
	Listeners   []string
}

var (
	listenerMap sync.Map // key: event name -> []base.Listener[T]
	eventInfos  sync.Map // key: event name -> EventInfo
)

// Register 注册监听
func Register(listener base.Listener, e base.Event) {
	name := e.Name()
	desc := e.Description()

	// 获取当前已注册监听
	var listen []base.Listener
	if v, ok := listenerMap.Load(name); ok {
		listen = v.([]base.Listener)
	}

	// 添加新的监听
	listen = append(listen, listener)
	listenerMap.Store(name, listen)

	// 记录事件信息
	if v, ok := eventInfos.Load(name); ok {
		info := v.(EventInfo)
		info.Listeners = append(info.Listeners, fmt.Sprintf("%T", listener))
		eventInfos.Store(name, info)
	} else {
		eventInfos.Store(name, EventInfo{
			Name:        name,
			Description: desc,
			Listeners:   []string{fmt.Sprintf("%T", listener)},
		})
	}
}

// Publish 发布事件
func Publish(ctx context.Context, e base.Event) {
	message.GetEventBus().Publish(debugger.TopicListener, debugger.ListenerEvent{
		TraceId:     ctx.Value(ctxkey.TraceIdKey).(string),
		Name:        e.Name(),
		Description: e.Description(),
		Data:        e,
	})

	if v, ok := listenerMap.Load(e.Name()); ok {
		for _, listener := range v.([]base.Listener) {
			// 可替换goroutine为队列
			go listener.Handle(e)
		}
	} else {
		color.Yellow("未找到事件监听: %s", e.Name())
	}
}

// EventList 已注册事件列表
func EventList() []EventInfo {
	var list []EventInfo
	eventInfos.Range(func(_, value any) bool {
		list = append(list, value.(EventInfo))
		return true
	})
	return list
}

// DebugPrint 打印所有注册事件信息
func DebugPrint() {
	color.Cyan("==== 当前已注册事件 ====")
	eventInfos.Range(func(_, value any) bool {
		info := value.(EventInfo)
		fmt.Printf("事件: %s\n描述: %s\n监听:\n", info.Name, info.Description)
		for _, l := range info.Listeners {
			fmt.Printf("  - %s\n", l)
		}
		fmt.Println("----------------------")
		return true
	})
}
