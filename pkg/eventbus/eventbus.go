package eventbus

import (
	"context"
	"fmt"
	"gin/common/base"
	"gin/common/ctxkey"
	"gin/pkg/debugger"
	"gin/pkg/message"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"sort"
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

	// 按名称排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].Name < events[j].Name
	})

	// 计算最大名称宽度和描述宽度(使用显示宽度)
	maxNameLen := 0
	maxDescLen := 0
	for _, info := range events {
		nameLen := runewidth.StringWidth(info.Name)
		if nameLen > maxNameLen {
			maxNameLen = nameLen
		}

		descLen := runewidth.StringWidth(info.Description)
		if descLen > maxDescLen {
			maxDescLen = descLen
		}
	}

	// 设置最小宽度
	if maxNameLen < 20 {
		maxNameLen = 20
	}
	if maxDescLen < 35 {
		maxDescLen = 35
	}

	// 计算标题的显示宽度
	titleNameLen := runewidth.StringWidth("事件名称")
	titleDescLen := runewidth.StringWidth("描述")
	if titleNameLen > maxNameLen {
		maxNameLen = titleNameLen
	}
	if titleDescLen > maxDescLen {
		maxDescLen = titleDescLen
	}

	// 计算总宽度
	totalWidth := maxNameLen + maxDescLen + 7

	// 打印顶部边框
	color.Yellow("┌" + strings.Repeat("─", totalWidth-2) + "┐")

	// 打印标题行
	titleLine := fmt.Sprintf("│ %s   %s "+color.YellowString("│"),
		color.HiWhiteString(padRight("事件名称", maxNameLen)),
		color.HiWhiteString(padRight("描述", maxDescLen)))
	color.Yellow(titleLine)

	// 打印分隔线
	color.Yellow("├" + strings.Repeat("─", maxNameLen+2) + "─" + strings.Repeat("─", maxDescLen+2) + "┤")

	var listeners int
	// 打印事件列表
	for _, info := range events {
		// 事件行
		coloredName := color.GreenString(padRight(info.Name, maxNameLen))
		coloredDesc := color.WhiteString(padRight(info.Description, maxDescLen))
		contentLine := fmt.Sprintf("│ %s   %s "+color.YellowString("│"), coloredName, coloredDesc)
		color.Yellow(contentLine)

		// 打印监听器
		for i, listener := range info.Listeners {
			prefix := "├─ "
			if i == len(info.Listeners)-1 {
				prefix = "└─ "
			}
			listeners++

			// 构建完整行使用固定格式
			// 名称区域：maxNameLen+2个空格+前缀+监听器名称
			nameArea := strings.Repeat(" ", maxNameLen+2) + prefix + listener
			// 确保总长度与事件行一致
			fullLine := fmt.Sprintf("│%-*s│", totalWidth-2, nameArea)
			color.Yellow(fullLine)
		}
	}

	// 打印底部边框
	color.Yellow("└" + strings.Repeat("─", totalWidth-2) + "┘")

	// 打印统计信息
	color.Cyan(fmt.Sprintf("总计 %d 个事件 %d 个监听\n", len(events), listeners))
}

// padRight 右侧填充空格,支持中文字符
func padRight(s string, width int) string {
	currentWidth := runewidth.StringWidth(s)
	if currentWidth >= width {
		return s
	}
	padding := width - currentWidth
	return s + strings.Repeat(" ", padding)
}
