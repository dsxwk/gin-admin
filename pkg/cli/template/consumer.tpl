package {{.Package}}

import (
	"context"
	"gin/app/facade"
	"gin/pkg"
	"gin/pkg/serviceprovider/queue"
	{{- if eq .Type "kafka"}}
	"time"

	"github.com/segmentio/kafka-go"
	{{- end}}
)

// {{.TypeTitle}}{{.Name}}Payload 消息体
type {{.TypeTitle}}{{.Name}}Payload struct {
}

// New{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer 创建消费者实例
func New{{.Name}}{{if .IsDelay}}Delay{{end}}Consumer() *queue.ConsumerAdapter[{{.TypeTitle}}{{.Name}}Payload] {
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
	{{- else if eq .Type "rabbitmq"}}
	log := facade.Log()
	mq, err := queue.NewRabbitMQ(facade.Config(), log, facade.Event().Bus())
	if err != nil {
		log.Error(pkg.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}
	{{- end}}

	var driver queue.ConsumerDriver
	{{- if eq .Type "kafka"}}
	driver = &queue.KafkaConsumer{
		Kafka: kfk,
		Topic: "{{.Topic}}",
		Group: "{{.Group}}",
	}
	{{- else if eq .Type "rabbitmq"}}
	driver = &queue.RabbitmqConsumer{
		Mq:       mq,
		Queue:    "{{.Queue}}",
		Exchange: "{{.Exchange}}",
		Routing:  "{{.Routing}}",
	}
	{{- else}}
	driver = &queue.RedisConsumer{
		Queue:     "{{.Queue}}",
		GetClient: facade.RedisClient,
		Log:       facade.Log(),
	}
	{{- end}}

	return queue.NewConsumerAdapter(
		queue.ConsumerConfig{
			Name:        "{{if .IsDelay}}{{.LowerName}}_delay{{else}}{{.LowerName}}{{end}}",
			Description: "{{.Description}}",
			Connection:  "{{.Type}}",
			Retry:       {{.Retry}},
			IsDelay:     {{.IsDelay}},
		},
		driver,
		handle{{.TypeTitle}}{{.Name}}{{if .IsDelay}}Delay{{end}},
	)
}

// handle{{.TypeTitle}}{{.Name}}{{if .IsDelay}}Delay{{end}} 处理{{.TypeTitle}}消息
func handle{{.TypeTitle}}{{.Name}}{{if .IsDelay}}Delay{{end}}(_ context.Context, data *{{.TypeTitle}}{{.Name}}Payload) error {
	facade.Log().Info(pkg.Sprintf("{{.TypeTitle}} Received Msg: %v", data))
	// 待处理业务逻辑
	return nil
}
