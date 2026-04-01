package {{.Package}}

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/config"
	"gin/pkg"
	"gin/pkg/logger"
	"gin/pkg/queue"
	{{- if eq .Type "kafka"}}
	"github.com/segmentio/kafka-go"
	"time"
	{{- end}}
)

// {{.Name}}{{if .IsDelay}}Delay{{end}}Consumer {{.Type}}消费者
type {{.Name}}{{if .IsDelay}}Delay{{end}}Consumer struct {
	{{- if eq .Type "kafka"}}
	*base.KafkaConsumer
	{{- else}}
	*base.RabbitmqConsumer
	{{- end}}
}

// New{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer 创建消费者实例
func New{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer() *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer {
	cfg := facade.Config.Get()
	log := facade.Log.Logger()
	bus := facade.Message.GetBus()

	{{- if eq .Type "kafka"}}
	kfk := base.NewKafka(cfg, log, bus)
	kfk.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Kafka.Brokers,
		Topic:          "{{.Topic}}",
		GroupID:        "{{.Group}}",
		MinBytes:       1,
		MaxBytes:       10e6,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 0,
		MaxWait:        5 * time.Second,
	})

	return &{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer{
		KafkaConsumer: &base.KafkaConsumer{
			Kafka:        kfk,
			Topic:        "{{.Topic}}",
			Group:        "{{.Group}}",
			Retry:        {{.Retry}},
			IsDelayQueue: {{.IsDelay}},
		},
	}
	{{- else}}
	mq, err := base.NewRabbitMQ(cfg, log, bus)
	if err != nil {
		log.Errorf("RabbitMQ连接失败: %v", err)
		return nil
	}

	return &{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer{
		RabbitmqConsumer: &base.RabbitmqConsumer{
			Mq:           mq,
			Queue:        "{{.Queue}}",
			Exchange:     "{{.Exchange}}",
			Routing:      "{{.Routing}}",
			IsDelayQueue: {{.IsDelay}},
			Retry:        {{.Retry}},
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

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) Start(cfg *config.Config, log *logger.Logger) error {
	{{- if eq .Type "kafka"}}
	c.KafkaConsumer.Start(c)
	{{- else}}
	c.RabbitmqConsumer.Start(c)
	{{- end}}
	log.Info(pkg.Sprintf("%s消费者启动成功: %s", "{{.TypeTitle}}", c.Name()))
	return nil
}

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) Stop() error {
	{{- if eq .Type "kafka"}}
	return c.KafkaConsumer.Stop()
	{{- else}}
	return c.RabbitmqConsumer.Stop()
	{{- end}}
}

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) Enabled(cfg *config.Config) bool {
	{{- if eq .Type "kafka"}}
	return cfg.Kafka.Enabled
	{{- else}}
	return cfg.Rabbitmq.Enabled
	{{- end}}
}

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) Status() queue.ConsumerStatus {
	{{- if eq .Type "kafka"}}
	return c.KafkaConsumer.Status()
	{{- else}}
	return c.RabbitmqConsumer.Status()
	{{- end}}
}

func (c *{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer) Handle(msg string) error {
	facade.Log.Info(pkg.Sprintf("%s Received Msg: %s", "{{.TypeTitle}}", msg))
	// todo 处理业务逻辑
	return nil
}

func init() {
	queue.GetConsumerRegistry().Register(New{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer())
}
