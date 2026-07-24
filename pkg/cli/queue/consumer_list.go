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

type ConsumerList struct{}

func (s *ConsumerList) Name() string               { return "consumer:list" }
func (s *ConsumerList) Description() string        { return "消费者列表" }
func (s *ConsumerList) Help() []base.CommandOption { return []base.CommandOption{} }

func (s *ConsumerList) Execute(values map[string]string) {
	consumers := queue.GetConsumerRegistry().GetAll()
	if len(consumers) == 0 {
		color.Yellow("no registered consumers")
		return
	}

	sort.Slice(consumers, func(i, j int) bool { return consumers[i].Name() < consumers[j].Name() })

	maxName, maxConn, maxDelay, maxDesc := 0, 0, 0, 0
	for _, c := range consumers {
		if w := runewidth.StringWidth(c.Name()); w > maxName {
			maxName = w
		}
		if w := runewidth.StringWidth(c.Connection()); w > maxConn {
			maxConn = w
		}
		s := fmt.Sprintf("%v", c.IsDelay())
		if w := runewidth.StringWidth(s); w > maxDelay {
			maxDelay = w
		}
		if w := runewidth.StringWidth(c.Description()); w > maxDesc {
			maxDesc = w
		}
	}

	tName := runewidth.StringWidth("消费者名称")
	tConn := runewidth.StringWidth("连接")
	tDelay := runewidth.StringWidth("延迟队列")
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
	if maxDesc < 20 {
		maxDesc = 20
	}

	tw := maxName + maxConn + maxDelay + maxDesc + 13

	color.Yellow("┌" + strings.Repeat("─", tw-2) + "┐")

	tl := fmt.Sprintf("│ %s   %s   %s   %s "+color.YellowString("│"),
		color.HiWhiteString(padRight("消费者名称", maxName)),
		color.HiWhiteString(padRight("连接", maxConn)),
		color.HiWhiteString(padRight("延迟队列", maxDelay)),
		color.HiWhiteString(padRight("描述", maxDesc)))
	color.Yellow(tl)

	color.Yellow("├" + strings.Repeat("─", tw-2) + "┤")

	for _, c := range consumers {
		ds := fmt.Sprintf("%v", c.IsDelay())
		cl := fmt.Sprintf("│ %s   %s   %s   %s "+color.YellowString("│"),
			color.GreenString(padRight(c.Name(), maxName)),
			color.YellowString(padRight(c.Connection(), maxConn)),
			color.WhiteString(padRight(ds, maxDelay)),
			color.WhiteString(padRight(c.Description(), maxDesc)))
		color.Yellow(cl)
	}

	color.Yellow("└" + strings.Repeat("─", tw-2) + "┘")
	color.Cyan(fmt.Sprintf("总计 %d 消费者", len(consumers)))
}

func init() { cli.Register(&ConsumerList{}) }
