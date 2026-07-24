package base

import (
	"context"
	"gin/common/ctxkey"
	"gin/common/flag"
	"gin/config"
	"gin/pkg"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/logger"
	"gin/pkg/serviceprovider/message"
	"gin/pkg/serviceprovider/queue"
	"github.com/goccy/go-json"
	"github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

// RabbitMQ RabbitMQ连接
type RabbitMQ struct {
	Conn    *amqp091.Connection
	Channel *amqp091.Channel
	Conf    *config.Config
	Log     *logger.Logger
	Message *message.Event
}

func NewRabbitMQ(conf *config.Config, log *logger.Logger, bus *message.Event) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(conf.Queue.Rabbitmq.Url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQ{Conn: conn, Channel: ch, Conf: conf, Log: log, Message: bus}, nil
}

func (r *RabbitMQ) Close() error {
	var err1, err2 error
	if r.Channel != nil {
		err1 = r.Channel.Close()
		r.Channel = nil
	}
	if r.Conn != nil {
		err2 = r.Conn.Close()
		r.Conn = nil
	}
	if err1 != nil {
		return err1
	}
	return err2
}

// RabbitmqConsumer RabbitMQ消费者基类
type RabbitmqConsumer struct {
	Mq       *RabbitMQ
	Queue    string
	Exchange string
	Routing  string
	status   queue.ConsumerStatus
	statusMu sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

func (c *RabbitmqConsumer) Status() queue.ConsumerStatus {
	c.statusMu.RLock()
	defer c.statusMu.RUnlock()
	return c.status
}

func (c *RabbitmqConsumer) setStatus(status queue.ConsumerStatus) {
	c.statusMu.Lock()
	defer c.statusMu.Unlock()
	c.status = status
}

func (c *RabbitmqConsumer) Start(h interface{}) {
	c.setStatus(queue.ConsumerStatusRunning)
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go func() {
		defer c.setStatus(queue.ConsumerStatusStopped)
		for {
			select {
			case <-c.ctx.Done():
				flag.Infof("[RabbitMq] 消费者 %s 已停止", c.Queue)
				return
			default:
				c.consumeLoop(h)
			}
		}
	}()
}

func (c *RabbitmqConsumer) consumeLoop(h interface{}) {
	if c.Mq == nil || c.Mq.Channel == nil {
		time.Sleep(time.Second)
		return
	}

	args := amqp091.Table{}
	exchangeType := "direct"
	if h.(queue.Consumer).IsDelay() {
		exchangeType = "x-delayed-message"
		args["x-delayed-type"] = "direct"
	}

	if err := c.Mq.Channel.ExchangeDeclare(c.Exchange, exchangeType, true, false, false, false, args); err != nil {
		c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] ExchangeDeclare error: %v", err))
		time.Sleep(time.Second)
		return
	}

	if _, err := c.Mq.Channel.QueueDeclare(c.Queue, true, false, false, false, nil); err != nil {
		c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] QueueDeclare error: %v", err))
		time.Sleep(time.Second)
		return
	}

	if err := c.Mq.Channel.QueueBind(c.Queue, c.Routing, c.Exchange, false, nil); err != nil {
		c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] QueueBind error: %v", err))
		time.Sleep(time.Second)
		return
	}

	msgs, err := c.Mq.Channel.Consume(c.Queue, "", false, false, false, false, nil)
	if err != nil {
		c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] Consume error: %v", err))
		time.Sleep(time.Second)
		return
	}

	for msg := range msgs {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.handleMessage(msg, h)
		}
	}
}

func (c *RabbitmqConsumer) handleMessage(msg amqp091.Delivery, h interface{}) {
	maxRetry := h.(queue.Consumer).Retry()
	retry := 0
	for {
		err := queue.TryHandle(h, msg.Body)
		if err == nil {
			if ackErr := msg.Ack(false); ackErr != nil {
				c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] Ack error: %v", ackErr))
			}
			break
		}

		retry++
		if retry >= maxRetry {
			c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] Retry failed: %s", string(msg.Body)))
			if ackErr := msg.Ack(false); ackErr != nil {
				c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] Ack error: %v", ackErr))
			}
			break
		}
		time.Sleep(time.Second)
	}
}

func (c *RabbitmqConsumer) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}
	if c.Mq != nil {
		return c.Mq.Close()
	}
	return nil
}

// RabbitmqProducer RabbitMQ生产者基类
type RabbitmqProducer struct {
	Mq       *RabbitMQ
	Queue    string
	Exchange string
	Routing  string
	Owner    queue.Producer
	Headers  amqp091.Table
}

func (p *RabbitmqProducer) newQueue() error {
	args := amqp091.Table{}
	exchangeType := "direct"
	if p.Owner != nil && p.Owner.IsDelay() {
		exchangeType = "x-delayed-message"
		args["x-delayed-type"] = "direct"
	}

	if err := p.Mq.Channel.ExchangeDeclare(p.Exchange, exchangeType, true, false, false, false, args); err != nil {
		return err
	}

	if _, err := p.Mq.Channel.QueueDeclare(p.Queue, true, false, false, false, nil); err != nil {
		return err
	}

	return p.Mq.Channel.QueueBind(p.Queue, p.Routing, p.Exchange, false, nil)
}

func (p *RabbitmqProducer) Publish(ctx context.Context, msg any) error {
	// 自动序列化为JSON
	var body []byte
	switch v := msg.(type) {
	case []byte:
		body = v
	case string:
		body = []byte(v)
	default:
		var err error
		body, err = json.Marshal(msg)
		if err != nil {
			return err
		}
	}
	start := time.Now()

	if err := p.newQueue(); err != nil {
		p.Mq.Log.Error(pkg.Sprintf("[RabbitMq] newQueue error: %v", err))
		return err
	}

	headers := p.Headers
	if headers == nil {
		headers = amqp091.Table{}
	}
	if p.Owner != nil && p.Owner.IsDelay() && p.Owner.DelayMs() > 0 {
		headers["x-delay"] = p.Owner.DelayMs()
	}

	pub := amqp091.Publishing{
		ContentType: "application/json",
		Body:        body,
		Headers:     headers,
	}

	err := p.Mq.Channel.Publish(p.Exchange, p.Routing, false, false, pub)
	if err != nil {
		p.Mq.Log.Error(pkg.Sprintf("[RabbitMq] Publish error: %v", err))
	}

	traceId := "unknown"
	if ctx != nil {
		if id := ctx.Value(ctxkey.TraceIdKey); id != nil {
			if s, ok := id.(string); ok && s != "" {
				traceId = s
			}
		}
	}

	p.Mq.Message.Publish(debugger.TopicMq, debugger.MqEvent{
		TraceId: traceId,
		Driver:  "rabbitmq",
		Topic:   p.Exchange + ":" + p.Routing,
		Message: string(body),
		Ms:      float64(time.Since(start).Milliseconds()),
		Extra: map[string]any{
			"exchange": p.Exchange,
			"routing":  p.Routing,
			"queue":    p.Queue,
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
