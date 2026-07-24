package make

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg/cli"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type MakeQueue struct {
	base.BaseCommand
}

// Name 命令名称
func (m *MakeQueue) Name() string {
	return "make:queue"
}

// Description 命令描述
func (m *MakeQueue) Description() string {
	return "创建消息队列(Kafka/RabbitMQ/Redis)"
}

// Help 命令参数
func (m *MakeQueue) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			Flag:     base.Flag{Short: "n", Long: "name"},
			Desc:     "队列名称, 如: order",
			Required: true,
		},
		{
			Flag:     base.Flag{Short: "c", Long: "connection"},
			Desc:     "连接类型: kafka, rabbitmq, redis(默认取queue.connection配置)",
			Required: false,
		},
		{
			Flag:     base.Flag{Short: "d", Long: "delay"},
			Desc:     "是否延迟队列: true/false",
			Required: false,
		},
		{
			Flag:     base.Flag{Short: "D", Long: "desc"},
			Desc:     "队列描述",
			Required: false,
		},
	}
}

// Execute 执行命令
func (m *MakeQueue) Execute(values map[string]string) {
	name := values["name"]
	conn := values["connection"]
	isDelay := m.StringToBool(values["delay"])

	if conn == "" {
		cfg := facade.Config()
		if cfg != nil {
			conn = cfg.Queue.Connection
		}
		if conn == "" {
			conn = "redis"
		}
	}

	if conn != "kafka" && conn != "rabbitmq" && conn != "redis" {
		flag.Errorf("不支持的连接类型: %s, 仅支持 kafka, rabbitmq, redis", conn)
		os.Exit(1)
	}

	m.generateQueue(conn, name, isDelay, values)
}

// generateQueue 生成队列文件
func (m *MakeQueue) generateQueue(conn, name string, isDelay bool, values map[string]string) {
	queueName := strings.ToLower(name)
	camelName := lo.PascalCase(name)
	lowerName := lo.CamelCase(name)

	desc := values["desc"]
	if desc == "" {
		desc = name
	}

	typeTitle := map[string]string{"kafka": "Kafka", "rabbitmq": "RabbitMQ", "redis": "Redis"}[conn]

	// 根据name自动生成Topic/Key/Group/Queue/Exchange/Routing
	topic := queueName
	key := queueName + "_key"
	group := queueName + "_group"
	queue := queueName
	exchange := queueName + "_exchange"
	routing := queueName

	data := map[string]interface{}{
		"Package":     "consumer",
		"Name":        camelName,
		"LowerName":   lowerName,
		"Type":        conn,
		"TypeTitle":   typeTitle,
		"IsDelay":     isDelay,
		"Retry":       3,
		"DelayMs":     0,
		"Description": desc,
		"Topic":       topic,
		"Key":         key,
		"Group":       group,
		"Queue":       queue,
		"Exchange":    exchange,
		"Routing":     routing,
	}

	consumerFile := m.GetTemplate("consumer")
	consumerTpl, err := template.ParseFiles(consumerFile)
	if err != nil {
		flag.Errorf("Error parsing consumer template: %s", err.Error())
		os.Exit(1)
	}
	producerFile := m.GetTemplate("producer")
	producerTpl, err := template.ParseFiles(producerFile)
	if err != nil {
		flag.Errorf("Error parsing producer template: %s", err.Error())
		os.Exit(1)
	}

	f1 := filepath.Join("app/queue", "consumer", queueName+".go")
	cf := m.CheckDirAndFile(f1)
	if cf == nil {
		return
	}

	f2 := filepath.Join("app/queue", "producer", queueName+".go")
	pf := m.CheckDirAndFile(f2)
	if pf == nil {
		return
	}

	data["Package"] = "consumer"
	err = consumerTpl.Execute(cf, data)
	if err != nil {
		flag.Errorf("Error executing consumer template: %s", err.Error())
		os.Exit(1)
	}

	data["Package"] = "producer"
	err = producerTpl.Execute(pf, data)
	if err != nil {
		flag.Errorf("Error executing producer template: %s", err.Error())
		os.Exit(1)
	}

	flag.Successf("消费者文件: %s 生成成功!", f1)
	flag.Successf("生产者文件: %s 生成成功!", f2)
}

func init() {
	cli.Register(&MakeQueue{})
}
