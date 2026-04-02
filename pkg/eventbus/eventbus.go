package eventbus

import (
	"context"
	"fmt"
	"gin/common/base"
	"gin/common/ctxkey"
	"gin/pkg/debugger"
	"gin/pkg/message"
	"github.com/fatih/color"
	"strings"
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
	// 收集所有事件
	var events []EventInfo
	eventInfos.Range(func(_, value any) bool {
		events = append(events, value.(EventInfo))
		return true
	})

	if len(events) == 0 {
		color.Yellow("暂无注册的事件")
		return
	}

	// 计算最大宽度
	maxNameLen := 0
	maxDescLen := 0
	for _, info := range events {
		if len(info.Name) > maxNameLen {
			maxNameLen = len(info.Name)
		}
		if len(info.Description) > maxDescLen {
			maxDescLen = len(info.Description)
		}
	}

	// 设置最小宽度
	if maxNameLen < 15 {
		maxNameLen = 15
	}
	if maxDescLen < 30 {
		maxDescLen = 30
	}

	totalWidth := maxNameLen + maxDescLen + 8

	// 打印标题
	color.Cyan("\n" + strings.Repeat("=", totalWidth))
	color.Cyan(fmt.Sprintf("%-*s %-*s", maxNameLen+2, "事件名称", maxDescLen+2, "描述"))
	color.Cyan(strings.Repeat("=", totalWidth))

	// 打印事件
	for _, info := range events {
		fmt.Printf("%s %s\n",
			color.GreenString(fmt.Sprintf("%-*s", maxNameLen+2, info.Name)),
			color.YellowString(fmt.Sprintf("%-*s", maxDescLen+2, info.Description)),
		)

		// 打印监听器(缩进显示)
		for i, listener := range info.Listeners {
			prefix := "  └─ "
			if i == len(info.Listeners)-1 {
				prefix = "  └─ "
			} else {
				prefix = "  ├─ "
			}
			fmt.Printf("%s%s\n",
				strings.Repeat(" ", maxNameLen+2),
				color.CyanString("%s%s", prefix, listener),
			)
		}
	}

	color.Cyan(strings.Repeat("=", totalWidth))
	color.Cyan("总计 %d 个事件\n", len(events))
}
