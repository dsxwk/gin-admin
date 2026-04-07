package make

import (
	"gin/common/base"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/cli"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type MakeQueue struct {
	base.BaseCommand
}

func (m *MakeQueue) Name() string {
	return "make:queue"
}

func (m *MakeQueue) Description() string {
	return "创建消息队列(Kafka/RabbitMQ)"
}

func (m *MakeQueue) Help() []base.CommandOption {
	return []base.CommandOption{
		{
			base.Flag{Short: "t", Long: "type"},
			"队列类型: kafka 或 rabbitmq",
			true,
		},
		{
			base.Flag{Short: "n", Long: "name"},
			"队列名称, 如: order",
			true,
		},
		{
			base.Flag{Short: "d", Long: "delay"},
			"是否延迟队列: true/false",
			false,
		},
		{
			base.Flag{Short: "D", Long: "desc"},
			"队列描述",
			false,
		},
		// Kafka参数
		{
			base.Flag{Short: "T", Long: "topic"},
			"Kafka主题",
			false,
		},
		{
			base.Flag{Short: "k", Long: "key"},
			"Kafka消息Key",
			false,
		},
		{
			base.Flag{Short: "g", Long: "group"},
			"Kafka消费组",
			false,
		},
		// RabbitMQ参数
		{
			base.Flag{Short: "q", Long: "queue"},
			"RabbitMQ队列名",
			false,
		},
		{
			base.Flag{Short: "e", Long: "exchange"},
			"RabbitMQ交换机",
			false,
		},
		{
			base.Flag{Short: "r", Long: "routing"},
			"RabbitMQ路由键",
			false,
		},
		// 通用参数
		{
			base.Flag{Short: "R", Long: "retry", Default: "3"},
			"重试次数, 默认3",
			false,
		},
		{
			base.Flag{Short: "m", Long: "delayMs", Default: "0"},
			"延迟毫秒数, 默认0",
			false,
		},
	}
}

func (m *MakeQueue) Execute(args []string) {
	values := m.ParseFlags(m.Name(), args, m.Help())

	queueType := values["type"]
	name := values["name"]
	isDelay := m.StringToBool(values["delay"])

	m.generateQueue(queueType, name, isDelay, values)
}

func init() {
	cli.Register(&MakeQueue{})
}

func (m *MakeQueue) generateQueue(queueType, name string, isDelay bool, values map[string]string) {
	// 设置默认参数
	queueName := strings.ToLower(name)
	camelName := pkg.ToUpperCamel(name)
	lowerName := pkg.ToLowerCamel(name)

	retry, _ := pkg.StringToInt[int](values["retry"])
	delayMs, _ := pkg.StringToInt[int64](values["delayMs"])

	// 构建数据
	data := map[string]interface{}{
		"Package":     "consumer",
		"Name":        camelName,
		"LowerName":   lowerName,
		"Type":        queueType,
		"TypeTitle":   map[string]string{"kafka": "Kafka", "rabbitmq": "RabbitMQ"}[queueType],
		"IsDelay":     isDelay,
		"Retry":       retry,
		"DelayMs":     delayMs,
		"Description": values["desc"],
	}

	// Kafka参数
	if queueType == "kafka" {
		topic := values["topic"]
		if topic == "" {
			topic = queueName
		}
		group := values["group"]
		if group == "" {
			group = queueName + "_group"
		}
		key := values["key"]
		if key == "" {
			key = queueName + "_key"
		}
		data["Topic"] = topic
		data["Group"] = group
		data["Key"] = key
	} else {
		queue := values["queue"]
		if queue == "" {
			queue = queueName
		}
		exchange := values["exchange"]
		if exchange == "" {
			exchange = queueName + "_exchange"
		}
		routing := values["routing"]
		if routing == "" {
			routing = queueName
		}
		data["Queue"] = queue
		data["Exchange"] = exchange
		data["Routing"] = routing
	}

	// 获取模板
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

	f1 := filepath.Join("app/queue", queueType, "consumer", queueName+".go")
	cf := m.CheckDirAndFile(f1)
	if cf == nil {
		return
	}

	f2 := filepath.Join("app/queue", queueType, "producer", queueName+".go")
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
