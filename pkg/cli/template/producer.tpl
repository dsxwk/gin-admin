package {{.Package}}

import (
	"gin/app/facade"
	{{- if eq .Type "rabbitmq"}}
	"gin/pkg"
	{{- end}}
	"gin/pkg/serviceprovider/queue"
	{{- if eq .Type "kafka"}}

	"github.com/segmentio/kafka-go"
	{{- end}}
)

// New{{.Name}}{{if .IsDelay}}Delay{{end}}Producer 创建生产者实例
func New{{.Name}}{{if .IsDelay}}Delay{{end}}Producer() *queue.ProducerAdapter {
	{{- if eq .Type "kafka"}}
	cfg := facade.Config()
	kfk := queue.NewKafka(cfg, facade.Log(), facade.Event().Bus())
	kfk.Writer = &kafka.Writer{
		Addr:         kafka.TCP(cfg.Queue.Kafka.Brokers...),
		Topic:        "{{.Topic}}",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}
	{{- else if eq .Type "rabbitmq"}}
	log := facade.Log()
	mq, err := queue.NewRabbitMQ(facade.Config(), log, facade.Event().Bus())
	if err != nil {
		log.Error(pkg.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}
	{{- end}}

	var driver queue.ProducerDriver
	{{- if eq .Type "kafka"}}
	driver = &queue.KafkaProducer{
		Kafka: kfk,
		Topic: "{{.Topic}}",
		Key:   "{{.Key}}",
	}
	{{- else if eq .Type "rabbitmq"}}
	driver = &queue.RabbitmqProducer{
		Mq:       mq,
		Queue:    "{{.Queue}}",
		Exchange: "{{.Exchange}}",
		Routing:  "{{.Routing}}",
	}
	{{- else}}
	driver = &queue.RedisProducer{
		Queue:     "{{.Queue}}",
		GetClient: facade.RedisClient,
	}
	{{- end}}

	return queue.NewProducerAdapter(
		queue.ProducerConfig{
			Name:        "{{if .IsDelay}}{{.LowerName}}_delay{{else}}{{.LowerName}}{{end}}",
			Description: "{{.Description}}",
			Connection:  "{{.Type}}",
			IsDelay:     {{.IsDelay}},
			DelayMs:     {{.DelayMs}},
		},
		driver,
	)
}
