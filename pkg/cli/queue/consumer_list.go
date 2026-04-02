package queue

import (
	"fmt"
	"gin/common/base"
	"gin/pkg/cli"
	"gin/pkg/queue"
	"github.com/fatih/color"
	"sort"
	"strings"
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

	// 按名称排序
	sort.Slice(consumers, func(i, j int) bool {
		return consumers[i].Name() < consumers[j].Name()
	})

	// 计算最大名称宽度
	maxNameLen := 0
	for _, c := range consumers {
		if len(c.Name()) > maxNameLen {
			maxNameLen = len(c.Name())
		}
	}
	if maxNameLen < 20 {
		maxNameLen = 20
	}

	totalWidth := maxNameLen + 4

	// 打印标题
	color.Cyan("\n" + strings.Repeat("=", totalWidth))
	color.Cyan(fmt.Sprintf("%-*s", maxNameLen+2, "消费者名称"))
	color.Cyan(strings.Repeat("=", totalWidth))

	// 打印消费者列表
	for i, c := range consumers {
		prefix := "├─ "
		if i == len(consumers)-1 {
			prefix = "└─ "
		}
		// 只用一个格式符
		fmt.Printf("%s\n", color.GreenString(fmt.Sprintf("%-*s", maxNameLen+2, prefix+c.Name())))
	}

	color.Cyan(strings.Repeat("=", totalWidth))
	color.Cyan("总计 %d 个消费者\n", len(consumers))
}

func init() {
	cli.Register(&ConsumerList{})
}
