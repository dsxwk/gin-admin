package base

import (
	"context"
	"gin/common/ctxkey"
	"gin/config"
	"gin/pkg"
	"gin/pkg/debugger"
	"gin/pkg/logger"
	"gin/pkg/message"
	"gin/pkg/queue"
	"github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

// RabbitMQ RabbitMQ连接
type RabbitMQ struct {
	Conn    *amqp091.Connection // AMQP连接
	Channel *amqp091.Channel    // AMQP通道
	Conf    *config.Config      // 配置
	Log     *logger.Logger      // 日志
	Message *message.EventBus   // 事件总线
}

// NewRabbitMQ 创建RabbitMQ连接
func NewRabbitMQ(conf *config.Config, log *logger.Logger, bus *message.EventBus) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(conf.Rabbitmq.Url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return &RabbitMQ{Conn: conn, Channel: ch, Conf: conf, Log: log, Message: bus}, nil
}

// Close 关闭RabbitMQ连接
func (r *RabbitMQ) Close() error {
	var err1, err2 error

	// 关闭通道
	if r.Channel != nil {
		err1 = r.Channel.Close()
		r.Channel = nil
	}

	// 关闭连接
	if r.Conn != nil {
		err2 = r.Conn.Close()
		r.Conn = nil
	}

	// 优先返回第一个错误
	if err1 != nil {
		return err1
	}
	return err2
}

// RabbitmqConsumer RabbitMQ消费者基类
// 封装了队列声明、消费循环、重试机制等通用逻辑
type RabbitmqConsumer struct {
	Mq           *RabbitMQ            // RabbitMQ连接
	Queue        string               // 队列名称
	Exchange     string               // 交换机名称
	Routing      string               // 路由键
	IsDelayQueue bool                 // 是否为延迟队列
	Retry        int                  // 重试次数
	status       queue.ConsumerStatus // 消费者状态
	statusMu     sync.RWMutex         // 状态锁
	ctx          context.Context      // 上下文
	cancel       context.CancelFunc   // 取消函数
}

// Status 获取消费者状态
func (c *RabbitmqConsumer) Status() queue.ConsumerStatus {
	c.statusMu.RLock()
	defer c.statusMu.RUnlock()
	return c.status
}

// setStatus 设置消费者状态
func (c *RabbitmqConsumer) setStatus(status queue.ConsumerStatus) {
	c.statusMu.Lock()
	defer c.statusMu.Unlock()
	c.status = status
}

// Start 启动消费者
// h: 消息处理实现Handle(msg string) error方法
func (c *RabbitmqConsumer) Start(h queue.Handler) {
	c.setStatus(queue.ConsumerStatusRunning)
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go func() {
		defer c.setStatus(queue.ConsumerStatusStopped)

		for {
			select {
			case <-c.ctx.Done():
				c.Mq.Log.Info(pkg.Sprintf("[RabbitMq] 消费者 %s 已停止", c.Queue))
				return
			default:
				c.consumeLoop(h)
			}
		}
	}()
}

// consumeLoop 消费循环
func (c *RabbitmqConsumer) consumeLoop(h queue.Handler) {
	// 检查连接
	if c.Mq == nil || c.Mq.Channel == nil {
		time.Sleep(time.Second)
		return
	}

	// 声明交换机
	args := amqp091.Table{}
	exchangeType := "direct"
	if c.IsDelayQueue {
		exchangeType = "x-delayed-message"
		args["x-delayed-type"] = "direct"
	}

	if err := c.Mq.Channel.ExchangeDeclare(
		c.Exchange,
		exchangeType,
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		args,
	); err != nil {
		c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] ExchangeDeclare error: %v", err))
		time.Sleep(time.Second)
		return
	}

	// 声明队列
	if _, err := c.Mq.Channel.QueueDeclare(
		c.Queue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	); err != nil {
		c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] QueueDeclare error: %v", err))
		time.Sleep(time.Second)
		return
	}

	// 绑定队列到交换机
	if err := c.Mq.Channel.QueueBind(
		c.Queue,
		c.Routing,
		c.Exchange,
		false, // noWait
		nil,   // args
	); err != nil {
		c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] QueueBind error: %v", err))
		time.Sleep(time.Second)
		return
	}

	// 开始消费
	msgs, err := c.Mq.Channel.Consume(
		c.Queue,
		"",    // consumer
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] Consume error: %v", err))
		time.Sleep(time.Second)
		return
	}

	// 处理消息
	for msg := range msgs {
		select {
		case <-c.ctx.Done():
			return
		default:
			c.handleMessage(msg, h)
		}
	}
}

// handleMessage 处理单条消息
func (c *RabbitmqConsumer) handleMessage(msg amqp091.Delivery, h queue.Handler) {
	retry := 0
	for {
		err := h.Handle(string(msg.Body))
		if err == nil {
			// 处理成功确认消息
			if ackErr := msg.Ack(false); ackErr != nil {
				c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] Ack error: %v", ackErr))
			}
			break
		}

		retry++
		if retry >= c.Retry {
			c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] Retry failed: %s", string(msg.Body)))
			// 重试次数用完确认消息(避免无限重试)
			if ackErr := msg.Ack(false); ackErr != nil {
				c.Mq.Log.Error(pkg.Sprintf("[RabbitMq] Ack error: %v", ackErr))
			}
			break
		}
		time.Sleep(time.Second)
	}
}

// Stop 停止消费者
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
// 封装了队列声明、消息发布、调试信息等通用逻辑
type RabbitmqProducer struct {
	Mq           *RabbitMQ     // RabbitMQ连接
	Queue        string        // 队列名称
	Exchange     string        // 交换机名称
	Routing      string        // 路由键
	IsDelayQueue bool          // 是否为延迟队列
	DelayMs      int64         // 延迟毫秒数
	Headers      amqp091.Table // 消息头
}

// newQueue 声明队列和交换机
func (p *RabbitmqProducer) newQueue() error {
	args := amqp091.Table{}
	exchangeType := "direct"
	if p.IsDelayQueue {
		exchangeType = "x-delayed-message"
		args["x-delayed-type"] = "direct"
	}

	// 声明交换机
	if err := p.Mq.Channel.ExchangeDeclare(
		p.Exchange,
		exchangeType,
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		args,
	); err != nil {
		return err
	}

	// 声明队列
	if _, err := p.Mq.Channel.QueueDeclare(
		p.Queue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	); err != nil {
		return err
	}

	// 绑定队列到交换机
	return p.Mq.Channel.QueueBind(
		p.Queue,
		p.Routing,
		p.Exchange,
		false, // noWait
		nil,   // args
	)
}

// Publish 发布消息到RabbitMQ
// ctx: 上下文,用于链路追踪
// msg: 消息内容(字节数组)
func (p *RabbitmqProducer) Publish(ctx context.Context, msg []byte) error {
	start := time.Now()

	// 声明队列和交换机
	if err := p.newQueue(); err != nil {
		p.Mq.Log.Error(pkg.Sprintf("[RabbitMq] newQueue error: %v", err))
		return err
	}

	// 设置消息头
	headers := p.Headers
	if headers == nil {
		headers = amqp091.Table{}
	}
	if p.IsDelayQueue && p.DelayMs > 0 {
		headers["x-delay"] = p.DelayMs
	}

	// 构建发布消息
	pub := amqp091.Publishing{
		ContentType: "application/json",
		Body:        msg,
		Headers:     headers,
	}

	// 发布消息
	err := p.Mq.Channel.Publish(
		p.Exchange,
		p.Routing,
		false, // mandatory
		false, // immediate
		pub,
	)
	if err != nil {
		p.Mq.Log.Error(pkg.Sprintf("[RabbitMq] Publish error: %v", err))
	}

	// 获取traceId
	traceId := "unknown"
	if ctx != nil {
		if id := ctx.Value(ctxkey.TraceIdKey); id != nil {
			if s, ok := id.(string); ok && s != "" {
				traceId = s
			}
		}
	}

	// 发布调试事件
	p.Mq.Message.Publish(debugger.TopicMq, debugger.MqEvent{
		TraceId: traceId,
		Driver:  "rabbitmq",
		Topic:   p.Exchange + ":" + p.Routing,
		Message: string(msg),
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

// Close 关闭RabbitMQ连接
func (p *RabbitmqProducer) Close() error {
	if p.Mq != nil {
		return p.Mq.Close()
	}
	return nil
}
