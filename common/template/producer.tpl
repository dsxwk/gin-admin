package {{.Package}}

import (
	"context"
	"gin/app/facade"
	"gin/pkg/serviceprovider/queue"
	{{- if eq .Type "rabbitmq"}}
	"gin/pkg"
	{{- end}}
	{{- if eq .Type "kafka"}}
	"github.com/segmentio/kafka-go"
	{{- end}}
)

// {{.Name}}{{if .IsDelay}}Delay{{end}}Producer {{.TypeTitle}}生产者
type {{.Name}}{{if .IsDelay}}Delay{{end}}Producer struct {
	{{- if eq .Type "kafka"}}
	*queue.KafkaProducer
	{{- else if eq .Type "rabbitmq"}}
	*queue.RabbitmqProducer
	{{- else}}
	*queue.RedisProducer
	{{- end}}
}

// New{{.Name}}{{if .IsDelay}}Delay{{end}}Producer 创建生产者实例
func New{{.Name}}{{if .IsDelay}}Delay{{end}}Producer() *{{.Name}}{{if .IsDelay}}Delay{{end}}Producer {
	{{- if eq .Type "kafka"}}
	cfg := facade.Config()
	kfk := queue.NewKafka(cfg, facade.Log(), facade.Event().Bus())
	kfk.Writer = &kafka.Writer{
		Addr:         kafka.TCP(cfg.Queue.Kafka.Brokers...),
		Topic:        "{{.Topic}}",
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
	}

	p := &{{.Name}}{{if .IsDelay}}Delay{{end}}Producer{
		KafkaProducer: &queue.KafkaProducer{
			Kafka: kfk,
			Topic: "{{.Topic}}",
			Key:   "{{.Key}}",
		},
	}
	p.KafkaProducer.Owner = p
	return p
	{{- else if eq .Type "rabbitmq"}}
	log := facade.Log()
	mq, err := queue.NewRabbitMQ(facade.Config(), log, facade.Event().Bus())
	if err != nil {
		log.Error(pkg.Sprintf("RabbitMQ连接失败: %v", err))
		return nil
	}

	p := &{{.Name}}{{if .IsDelay}}Delay{{end}}Producer{
		RabbitmqProducer: &queue.RabbitmqProducer{
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
		RedisProducer: &queue.RedisProducer{
			Queue:     "{{.Queue}}",
			GetClient: facade.RedisClient,
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
	queue.GetProducerRegistry().RegisterFactory(func() queue.Producer {
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
		return New{{.Name}}{{if .IsDelay}}Delay{{end}}Producer()
	})
}
