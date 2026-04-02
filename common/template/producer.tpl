package {{.Package}}

import (
	"context"
	"gin/app/facade"
	"gin/common/base"
	"gin/pkg/queue"
	{{- if eq .Type "kafka"}}
	"github.com/segmentio/kafka-go"
	{{- end}}
	"sync"
)

// {{.Name}}{{if .IsDelay}}Delay{{end}}Producer {{.Type}}生产者
type {{.Name}}{{if .IsDelay}}Delay{{end}}Producer struct {
	{{- if eq .Type "kafka"}}
	*base.KafkaProducer
	{{- else}}
	*base.RabbitmqProducer
	{{- end}}
	initOnce sync.Once
}

// New{{.Name}}{{if .IsDelay}}Delay{{end}}Producer 创建生产者实例
func New{{.Name}}{{if .IsDelay}}Delay{{end}}Producer() *{{.Name}}{{if .IsDelay}}Delay{{end}}Producer {
	cfg := facade.Config.Get()
	log := facade.Log.Logger()
	bus := facade.Message.GetBus()

	{{- if eq .Type "kafka"}}
	kfk := base.NewKafka(cfg, log, bus)
	kfk.Writer = &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Brokers...),
		Topic:        "{{.Topic}}",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	return &{{.Name}}{{if .IsDelay}}Delay{{end}}Producer{
		KafkaProducer: &base.KafkaProducer{
			Kafka:        kfk,
			Topic:        "{{.Topic}}",
			Key:          "{{.Key}}",
			IsDelayQueue: {{.IsDelay}},
			DelayMs:      {{.DelayMs}},
		},
	}
	{{- else}}
	mq, err := base.NewRabbitMQ(cfg, log, bus)
	if err != nil {
		log.Errorf("RabbitMQ连接失败: %v", err)
		return nil
	}

	return &{{.Name}}{{if .IsDelay}}Delay{{end}}Producer{
		RabbitmqProducer: &base.RabbitmqProducer{
			Mq:           mq,
			Queue:        "{{.Queue}}",
			Exchange:     "{{.Exchange}}",
			Routing:      "{{.Routing}}",
			IsDelayQueue: {{.IsDelay}},
			DelayMs:      {{.DelayMs}},
		},
	}
	{{- end}}
}

func (p *{{.Name}}{{if .IsDelay}}Delay{{end}}Producer) Name() string {
	{{- if .IsDelay}}
	return "{{.LowerName}}_delay"
	{{- else}}
	return "{{.LowerName}}"
	{{- end}}
}

func (p *{{.Name}}{{if .IsDelay}}Delay{{end}}Producer) Publish(ctx context.Context, msg []byte) error {
	p.initOnce.Do(func() {})
	{{- if eq .Type "kafka"}}
	return p.KafkaProducer.Publish(ctx, msg)
	{{- else}}
	return p.RabbitmqProducer.Publish(ctx, msg)
	{{- end}}
}

func (p *{{.Name}}{{if .IsDelay}}Delay{{end}}Producer) Close() error {
	{{- if eq .Type "kafka"}}
	return p.KafkaProducer.Close()
	{{- else}}
	return p.RabbitmqProducer.Close()
	{{- end}}
}

func init() {
    cfg := facade.Config.Get()
	if cfg != nil && cfg.{{if eq .Type "kafka"}}Kafka{{else}}Rabbitmq{{end}}.Enabled {
	    queue.GetProducerRegistry().Register(New{{.Name}}{{if .IsDelay}}Delay{{end}}Producer())
	}
}
