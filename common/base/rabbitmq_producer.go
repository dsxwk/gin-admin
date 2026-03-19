package base

import (
	"context"
	"gin/common/ctxkey"
	"gin/config"
	"gin/pkg/debugger"
	"gin/pkg/logger"
	"gin/pkg/message"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

type RabbitMq struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
}

func NewAmqp(url string) (*RabbitMq, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMq{Conn: conn, Channel: ch}, nil
}

func (r *RabbitMq) Close() error {
	var err1, err2 error
	if r.Channel != nil {
		err1 = r.Channel.Close()
	}
	if r.Conn != nil {
		err2 = r.Conn.Close()
	}
	if err1 != nil {
		return err1
	}
	return err2
}

type RabbitmqProducer struct {
	Mq           *RabbitMq
	Queue        string
	Exchange     string
	Routing      string
	IsDelayQueue bool
	DelayMs      int64
	Headers      amqp091.Table
}

func NewRabbitMq() *RabbitMq {
	rmq, err := NewAmqp(config.NewConfig().Rabbitmq.Url)
	if err != nil {
		logger.NewLogger().Error("RabbitMq连接失败: " + err.Error())
	}
	return rmq
}

func (p *RabbitmqProducer) newQueue() error {
	args := amqp091.Table{}
	exchangeType := "direct"
	if p.IsDelayQueue {
		exchangeType = "x-delayed-message"
		args["x-delayed-type"] = "direct"
	}

	if err := p.Mq.Channel.ExchangeDeclare(
		p.Exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		args,
	); err != nil {
		return err
	}

	if _, err := p.Mq.Channel.QueueDeclare(p.Queue, true, false, false, false, nil); err != nil {
		return err
	}

	return p.Mq.Channel.QueueBind(p.Queue, p.Routing, p.Exchange, false, nil)
}

func (p *RabbitmqProducer) Publish(ctx context.Context, msg []byte) error {
	start := time.Now()
	if err := p.newQueue(); err != nil {
		return err
	}

	headers := p.Headers
	if headers == nil {
		headers = amqp091.Table{}
	}
	if p.IsDelayQueue && p.DelayMs > 0 {
		headers["x-delay"] = p.DelayMs
	}

	pub := amqp091.Publishing{
		ContentType: "application/json",
		Body:        msg,
		Headers:     headers,
	}

	err := p.Mq.Channel.Publish(p.Exchange, p.Routing, false, false, pub)

	message.GetEventBus().Publish(debugger.TopicMq, debugger.MqEvent{
		TraceId: ctx.Value(ctxkey.TraceIdKey).(string),
		Driver:  "rabbitmq",
		Topic:   p.Exchange + ":" + p.Routing,
		Message: string(msg),
		Key:     "",
		Group:   "",
		Ms:      float64(time.Since(start).Milliseconds()),
		Extra: map[string]any{
			"exchange": p.Exchange,
			"routing":  p.Routing,
			"err":      err,
		},
	})

	return err
}

func (p *RabbitmqProducer) Close() error {
	if p.Mq != nil {
		return p.Mq.Close()
	}
	return nil
}
