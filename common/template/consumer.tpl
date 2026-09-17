package {{.Package}}

import (
	"gin/app/facade"
	"gin/common/flag"
	"gin/config"
	"gin/pkg"
	"gin/pkg/serviceprovider/queue"
	{{- if eq .Type "kafka"}}
	"github.com/segmentio/kafka-go"
	"time"
	{{- end}}
)

// {{.Name}}{{if .IsDelay}}Delay{{end}}Consumer {{.TypeTitle}}消费者
type {{.Name}}{{if .IsDelay}}Delay{{end}}Consumer struct {
	{{- if eq .Type "kafka"}}
	*queue.KafkaConsumer
	{{- else if eq .Type "rabbitmq"}}
	*queue.RabbitmqConsumer
	{{- else}}
	*queue.RedisConsumer
	{{- end}}
}

// {{.TypeTitle}}{{.Name}}Payload 消息体
type {{.TypeTitle}}{{.Name}}Payload struct {
}

// New{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer 创建消费者实例
func New{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer() *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer {
	{{- if eq .Type "kafka"}}
	cfg := facade.Config()
	kfk := queue.NewKafka(cfg, facade.Log(), facade.Event().Bus())
	kfk.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Queue.Kafka.Brokers,
		Topic:          "{{.Topic}}",
		GroupID:        "{{.Group}}",
		MinBytes:       1,
		MaxBytes:       10e6,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 0,
		MaxWait:        5 * time.Second,
	})

	return &{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer{
		KafkaConsumer: &queue.KafkaConsumer{
			Kafka: kfk,
			Topic: "{{.Topic}}",
			Group: "{{.Group}}",
		},
	}
	{{- else if eq .Type "rabbitmq"}}
	log := facade.Log()
	mq, err := queue.NewRabbitMQ(facade.Config(), log, facade.Event().Bus())
	if err != nil {
		log.Error(pkg.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}

	return &{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer{
		RabbitmqConsumer: &queue.RabbitmqConsumer{
			Mq:      mq,
			Queue:   "{{.Queue}}",
			Exchange: "{{.Exchange}}",
			Routing: "{{.Routing}}",
		},
	}
	{{- else}}
	return &{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer{
		RedisConsumer: &queue.RedisConsumer{
			Queue:     "{{.Queue}}",
			GetClient: facade.RedisClient,
			Log:       facade.Log(),
		},
	}
	{{- end}}
}

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) Name() string {
	{{- if .IsDelay}}
	return "{{.LowerName}}_delay"
	{{- else}}
	return "{{.LowerName}}"
	{{- end}}
}

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) Description() string {
	return "{{.Description}}"
}

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) Connection() string {
	return "{{.Type}}"
}

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) Retry() int { return {{.Retry}} }

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) IsDelay() bool { return {{.IsDelay}} }

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) Start() error {
	{{- if eq .Type "kafka"}}
	c.KafkaConsumer.Start(c)
	{{- else if eq .Type "rabbitmq"}}
	c.RabbitmqConsumer.Start(c)
	{{- else}}
	c.RedisConsumer.Start(c)
	{{- end}}
	flag.Infof("{{.TypeTitle}}消费者启动成功: %s", c.Name())
	return nil
}

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) Stop() error {
	{{- if eq .Type "kafka"}}
	return c.KafkaConsumer.Stop()
	{{- else if eq .Type "rabbitmq"}}
	return c.RabbitmqConsumer.Stop()
	{{- else}}
	return c.RedisConsumer.Stop()
	{{- end}}
}

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) Enabled(cfg *config.Config) bool {
	{{- if eq .Type "kafka"}}
	return cfg.Queue.Kafka.Enabled
	{{- else if eq .Type "rabbitmq"}}
	return cfg.Queue.Rabbitmq.Enabled
	{{- else}}
	return true
	{{- end}}
}

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) NewPayload() any {
	return &{{.TypeTitle}}{{.Name}}Payload{}
}

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) Handle(payload any) error {
	data := payload.(*{{.TypeTitle}}{{.Name}}Payload)
	facade.Log().Info(pkg.Sprintf("{{.TypeTitle}} Received Msg: %v", data))
	// todo 处理业务逻辑
	return nil
}

func init() {
	queue.GetConsumerRegistry().RegisterFactory(func() queue.Consumer {
		{{- if eq .Type "kafka"}}
		cfg := facade.Config()
		if cfg == nil || !cfg.Queue.Kafka.Enabled {
			return nil
		}
		{{- else if eq .Type "rabbitmq"}}
		cfg := facade.Config()
		if cfg == nil || !cfg.Queue.Rabbitmq.Enabled {
			return nil
		}
		{{- end}}
		return New{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer()
	})
}
