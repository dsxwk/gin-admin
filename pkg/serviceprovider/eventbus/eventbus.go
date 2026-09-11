package eventbus

import (
	"context"
	"fmt"
	"gin/common/ctxkey"
	"gin/common/flag"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
)

const (
	// TopicListener 监听调试主题
	TopicListener = "debug:listener"
)

// ListenerEvent 监听调试事件
type ListenerEvent struct {
	TraceId     string
	Name        string
	Description string
	Data        any
}

type EventInfo struct {
	Name        string
	Description string
	Listeners   []string
}

// Registry 事件注册表
type Registry struct {
	mu        sync.RWMutex
	listeners map[string][]Listener[Event] // key: event name -> []Listener
	infos     map[string]*EventInfo        // key: event name -> EventInfo
}

// 默认注册表
var defaultRegistry = &Registry{
	listeners: make(map[string][]Listener[Event]),
	infos:     make(map[string]*EventInfo),
}

// listenerWrapper 泛型监听器包装器
// 用于将 Listener[T] 转换为 Listener[Event]
type listenerWrapper[T Event] struct {
	inner Listener[T]
}

func (w *listenerWrapper[T]) Handle(e Event) {
	if t, ok := e.(T); ok {
		w.inner.Handle(t)
	}
}

// Register 注册监听
func Register[T Event](listener Listener[T], event T) {
	defaultRegistry.mu.Lock()
	defer defaultRegistry.mu.Unlock()

	name := event.Name()

	// 包装监听器（Listener[T] -> Listener[Event]）
	wrapped := &listenerWrapper[T]{inner: listener}
	defaultRegistry.listeners[name] = append(defaultRegistry.listeners[name], wrapped)

	// 更新事件信息
	if info, ok := defaultRegistry.infos[name]; ok {
		info.Listeners = append(info.Listeners, fmt.Sprintf("%T", listener))
	} else {
		defaultRegistry.infos[name] = &EventInfo{
			Name:        name,
			Description: event.Description(),
			Listeners:   []string{fmt.Sprintf("%T", listener)},
		}
	}
}

// Publish 发布事件
func Publish[T Event](ctx context.Context, e T) {
	NewBus().Publish(TopicListener, ListenerEvent{
		TraceId:     ctxkey.GetTraceId(ctx),
		Name:        e.Name(),
		Description: e.Description(),
		Data:        e,
	})

	// 获取监听器
	defaultRegistry.mu.RLock()
	listeners := defaultRegistry.listeners[e.Name()]
	defaultRegistry.mu.RUnlock()

	if len(listeners) == 0 {
		flag.Warningf("未找到事件监听: %s", e.Name())
		return
	}

	// 异步执行监听器
	for _, listener := range listeners {
		// 可替换goroutine为队列
		go listener.Handle(e)
	}
}

// EventList 已注册事件列表
func EventList() []EventInfo {
	defaultRegistry.mu.RLock()
	defer defaultRegistry.mu.RUnlock()

	list := make([]EventInfo, 0, len(defaultRegistry.infos))
	for _, info := range defaultRegistry.infos {
		list = append(list, *info)
	}
	return list
}

// DebugPrint 打印所有注册事件信息
func DebugPrint() {
	// 获取所有事件
	defaultRegistry.mu.RLock()
	events := make([]EventInfo, 0, len(defaultRegistry.infos))
	for _, info := range defaultRegistry.infos {
		events = append(events, *info)
	}
	defaultRegistry.mu.RUnlock()

	if len(events) == 0 {
		flag.Warningf("暂无注册的事件")
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
