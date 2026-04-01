package queue

import (
	"gin/common/base"
	"gin/pkg/cli"
	"gin/pkg/queue"
	"github.com/fatih/color"
	"sort"
)

type ConsumerList struct{}

func (s *ConsumerList) Name() string {
	return "consumer:list"
}

func (s *ConsumerList) Description() string {
	return "消费者列表"
}

func (s *ConsumerList) Help() []base.CommandOption {
	return []base.CommandOption{}
}

func (s *ConsumerList) Execute(args []string) {
	consumers := queue.GetConsumerRegistry().GetAll()
	if len(consumers) == 0 {
		color.Yellow("暂无注册的消费者")
		return
	}

	// 提取名称并排序
	names := make([]string, 0, len(consumers))
	for _, c := range consumers {
		names = append(names, c.Name())
	}
	sort.Strings(names)

	// 打印表头
	color.Green("----------------------------------------")
	color.Green("消费者列表 (%d)\n", len(names))
	color.Green("----------------------------------------")

	for _, name := range names {
		color.Green("  %s", name)
	}

	color.Green("----------------------------------------")
}

func init() {
	cli.Register(&ConsumerList{})
}
