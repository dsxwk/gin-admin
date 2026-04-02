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

	// 按名称排序
	sort.Slice(producers, func(i, j int) bool {
		return producers[i].Name() < producers[j].Name()
	})

	// 计算最大名称宽度
	maxNameLen := 0
	for _, p := range producers {
		if len(p.Name()) > maxNameLen {
			maxNameLen = len(p.Name())
		}
	}
	if maxNameLen < 20 {
		maxNameLen = 20
	}

	totalWidth := maxNameLen + 4

	// 打印标题
	color.Cyan("\n" + strings.Repeat("=", totalWidth))
	color.Cyan(fmt.Sprintf("%-*s", maxNameLen+2, "生产者名称"))
	color.Cyan(strings.Repeat("=", totalWidth))

	// 打印生产者列表
	for i, p := range producers {
		prefix := "├─ "
		if i == len(producers)-1 {
			prefix = "└─ "
		}
		// 只用一个格式符
		fmt.Printf("%s\n", color.GreenString(fmt.Sprintf("%-*s", maxNameLen+2, prefix+p.Name())))
	}

	color.Cyan(strings.Repeat("=", totalWidth))
	color.Cyan("总计 %d 个生产者\n", len(producers))
}

func init() {
	cli.Register(&ProducerList{})
}
