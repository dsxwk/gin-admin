package {{.Package}}

import (
	"context"
	{{- if ne .Type "redis"}}
	"gin/app/facade"
	{{- end}}
	"gin/common/base"
	{{- if eq .Type "rabbitmq"}}
	"gin/pkg"
	{{- end}}
	"gin/pkg/serviceprovider/queue"
	{{- if eq .Type "kafka"}}
	"github.com/segmentio/kafka-go"
	{{- end}}
)

// {{.Name}}{{if .IsDelay}}Delay{{end}}Producer {{.TypeTitle}}生产者
type {{.Name}}{{if .IsDelay}}Delay{{end}}Producer struct {
	{{- if eq .Type "kafka"}}
	*base.KafkaProducer
	{{- else if eq .Type "rabbitmq"}}
	*base.RabbitmqProducer
	{{- else}}
	*base.RedisProducer
	{{- end}}
}

// New{{.Name}}{{if .IsDelay}}Delay{{end}}Producer 创建生产者实例
func New{{.Name}}{{if .IsDelay}}Delay{{end}}Producer() *{{.Name}}{{if .IsDelay}}Delay{{end}}Producer {
	{{- if eq .Type "kafka"}}
	cfg := facade.Config()
	kfk := base.NewKafka(cfg, facade.Log(), facade.Message())
	kfk.Writer = &kafka.Writer{
		Addr:         kafka.TCP(cfg.Queue.Kafka.Brokers...),
		Topic:        "{{.Topic}}",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	p := &{{.Name}}{{if .IsDelay}}Delay{{end}}Producer{
		KafkaProducer: &base.KafkaProducer{
			Kafka: kfk,
			Topic: "{{.Topic}}",
			Key:   "{{.Key}}",
		},
	}
	p.KafkaProducer.Owner = p
	return p
	{{- else if eq .Type "rabbitmq"}}
	log := facade.Log()
	mq, err := base.NewRabbitMQ(facade.Config(), log, facade.Message())
	if err != nil {
		log.Error(pkg.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}

	p := &{{.Name}}{{if .IsDelay}}Delay{{end}}Producer{
		RabbitmqProducer: &base.RabbitmqProducer{
			Mq:      mq,
			Queue:   "{{.Queue}}",
			Exchange: "{{.Exchange}}",
			Routing: "{{.Routing}}",
		},
	}
	p.RabbitmqProducer.Owner = p
	return p
	{{- else}}
	p := &{{.Name}}{{if .IsDelay}}Delay{{end}}Producer{
		RedisProducer: &base.RedisProducer{
			Queue: "{{.Queue}}",
		},
	}
	p.RedisProducer.Owner = p
	return p
	{{- end}}
}

func (p *{{.Name}}{{if .IsDelay}}Delay{{end}}Producer) Name() string {
	{{- if .IsDelay}}
	return "{{.LowerName}}_delay"
	{{- else}}
	return "{{.LowerName}}"
	{{- end}}
}

func (p *{{.Name}}{{if .IsDelay}}Delay{{end}}Producer) Description() string {
	return "{{.Description}}"
}

func (p *{{.Name}}{{if .IsDelay}}Delay{{end}}Producer) Connection() string {
	return "{{.Type}}"
}

func (p *{{.Name}}{{if .IsDelay}}Delay{{end}}Producer) IsDelay() bool { return {{.IsDelay}} }

func (p *{{.Name}}{{if .IsDelay}}Delay{{end}}Producer) DelayMs() int64 { return {{.DelayMs}} }

func (p *{{.Name}}{{if .IsDelay}}Delay{{end}}Producer) Publish(ctx context.Context, msg any) error {
	{{- if eq .Type "kafka"}}
	return p.KafkaProducer.Publish(ctx, msg)
	{{- else if eq .Type "rabbitmq"}}
	return p.RabbitmqProducer.Publish(ctx, msg)
	{{- else}}
	return p.RedisProducer.Publish(ctx, msg)
	{{- end}}
}

func (p *{{.Name}}{{if .IsDelay}}Delay{{end}}Producer) Close() error {
	{{- if eq .Type "kafka"}}
	return p.KafkaProducer.Close()
	{{- else if eq .Type "rabbitmq"}}
	return p.RabbitmqProducer.Close()
	{{- else}}
	return p.RedisProducer.Close()
	{{- end}}
}

func init() {
	{{- if eq .Type "kafka"}}
	cfg := facade.Config()
	if cfg != nil && cfg.Queue.Kafka.Enabled {
		if p := New{{.Name}}{{if .IsDelay}}Delay{{end}}Producer(); p != nil {
			queue.GetProducerRegistry().Register(p)
		}
	}
	{{- else if eq .Type "rabbitmq"}}
	cfg := facade.Config()
	if cfg != nil && cfg.Queue.Rabbitmq.Enabled {
		if p := New{{.Name}}{{if .IsDelay}}Delay{{end}}Producer(); p != nil {
			queue.GetProducerRegistry().Register(p)
		}
	}
	{{- else}}
	queue.GetProducerRegistry().Register(New{{.Name}}{{if .IsDelay}}Delay{{end}}Producer())
	{{- end}}
}