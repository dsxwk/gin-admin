package queue

import (
	"fmt"
	"gin/common/base"
	"gin/pkg/cli"
	"gin/pkg/serviceprovider/queue"
	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"
	"sort"
	"strings"
)

type ProducerList struct{}

func (s *ProducerList) Name() string               { return "producer:list" }
func (s *ProducerList) Description() string        { return "生产者列表" }
func (s *ProducerList) Help() []base.CommandOption { return []base.CommandOption{} }

func (s *ProducerList) Execute(values map[string]string) {
	producers := queue.GetProducerRegistry().GetAll()
	if len(producers) == 0 {
		color.Yellow("no registered producers")
		return
	}

	sort.Slice(producers, func(i, j int) bool { return producers[i].Name() < producers[j].Name() })

	maxName, maxConn, maxDelay, maxDelayMs, maxDesc := 0, 0, 0, 0, 0
	for _, p := range producers {
		if w := runewidth.StringWidth(p.Name()); w > maxName {
			maxName = w
		}
		if w := runewidth.StringWidth(p.Connection()); w > maxConn {
			maxConn = w
		}
		s := fmt.Sprintf("%v", p.IsDelay())
		if w := runewidth.StringWidth(s); w > maxDelay {
			maxDelay = w
		}
		s2 := fmt.Sprintf("%dms", p.DelayMs())
		if w := runewidth.StringWidth(s2); w > maxDelayMs {
			maxDelayMs = w
		}
		if w := runewidth.StringWidth(p.Description()); w > maxDesc {
			maxDesc = w
		}
	}

	tName := runewidth.StringWidth("生产者名称")
	tConn := runewidth.StringWidth("连接")
	tDelay := runewidth.StringWidth("延迟队列")
	tDelayMs := runewidth.StringWidth("延迟(ms)")
	tDesc := runewidth.StringWidth("描述")
	if tName > maxName {
		maxName = tName
	}
	if tConn > maxConn {
		maxConn = tConn
	}
	if tDelay > maxDelay {
		maxDelay = tDelay
	}
	if tDelayMs > maxDelayMs {
		maxDelayMs = tDelayMs
	}
	if tDesc > maxDesc {
		maxDesc = tDesc
	}
	if maxName < 20 {
		maxName = 20
	}
	if maxConn < 6 {
		maxConn = 6
	}
	if maxDelay < 8 {
		maxDelay = 8
	}
	if maxDelayMs < 8 {
		maxDelayMs = 8
	}
	if maxDesc < 20 {
		maxDesc = 20
	}

	tw := maxName + maxConn + maxDelay + maxDelayMs + maxDesc + 16

	color.Yellow("┌" + strings.Repeat("─", tw-2) + "┐")

	tl := fmt.Sprintf("│ %s   %s   %s   %s   %s "+color.YellowString("│"),
		color.HiWhiteString(padRight("生产者名称", maxName)),
		color.HiWhiteString(padRight("连接", maxConn)),
		color.HiWhiteString(padRight("延迟队列", maxDelay)),
		color.HiWhiteString(padRight("延迟(ms)", maxDelayMs)),
		color.HiWhiteString(padRight("描述", maxDesc)))
	color.Yellow(tl)

	color.Yellow("├" + strings.Repeat("─", tw-2) + "┤")

	for _, p := range producers {
		ds := fmt.Sprintf("%v", p.IsDelay())
		dms := fmt.Sprintf("%dms", p.DelayMs())
		cl := fmt.Sprintf("│ %s   %s   %s   %s   %s "+color.YellowString("│"),
			color.GreenString(padRight(p.Name(), maxName)),
			color.YellowString(padRight(p.Connection(), maxConn)),
			color.WhiteString(padRight(ds, maxDelay)),
			color.WhiteString(padRight(dms, maxDelayMs)),
			color.WhiteString(padRight(p.Description(), maxDesc)))
		color.Yellow(cl)
	}

	color.Yellow("└" + strings.Repeat("─", tw-2) + "┘")
	color.Cyan(fmt.Sprintf("总计 %d 生产者", len(producers)))
}

func padRight(s string, width int) string {
	currentWidth := runewidth.StringWidth(s)
	if currentWidth >= width {
		return s
	}
	padding := width - currentWidth
	return s + strings.Repeat(" ", padding)
}

func init() { cli.Register(&ProducerList{}) }
