package queue

import (
	"gin/common/base"
	"gin/pkg/cli"
	"gin/pkg/queue"
	"github.com/fatih/color"
	"sort"
)

type ProducerList struct{}

func (s *ProducerList) Name() string {
	return "producer:list"
}

func (s *ProducerList) Description() string {
	return "生产者列表"
}

func (s *ProducerList) Help() []base.CommandOption {
	return []base.CommandOption{}
}

func (s *ProducerList) Execute(args []string) {
	producers := queue.GetProducerRegistry().GetAll()

	if len(producers) == 0 {
		color.Yellow("暂无注册的生产者")
		return
	}

	// 提取名称并排序
	names := make([]string, 0, len(producers))
	for _, p := range producers {
		names = append(names, p.Name())
	}
	sort.Strings(names)

	// 打印表头
	color.Green("----------------------------------------")
	color.Green("生产者列表 (%d)\n", len(names))
	color.Green("----------------------------------------")

	for _, name := range names {
		color.Green("  %s", name)
	}

	color.Green("----------------------------------------")
}

func init() {
	cli.Register(&ProducerList{})
}
