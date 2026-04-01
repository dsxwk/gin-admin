package base

import (
	"context"
	"encoding/json"
	"errors"
	"gin/common/ctxkey"
	"gin/config"
	"gin/pkg"
	"gin/pkg/debugger"
	"gin/pkg/logger"
	"gin/pkg/message"
	"gin/pkg/queue"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

// Kafka Kafka连接
type Kafka struct {
	Writer  *kafka.Writer     // Kafka Writer
	Reader  *kafka.Reader     // Kafka Reader
	Conf    *config.Config    // 配置
	Log     *logger.Logger    // 日志
	Message *message.EventBus // 事件总线
}

// NewKafka 创建Kafka连接
func NewKafka(conf *config.Config, log *logger.Logger, bus *message.EventBus) *Kafka {
	return &Kafka{
		Conf:    conf,
		Log:     log,
		Message: bus,
	}
}

// KafkaConsumer Kafka消费者
// 封装了消息读取、延迟消息处理、重试机制等通用逻辑
type KafkaConsumer struct {
	Kafka        *Kafka               // Kafka连接
	Topic        string               // 主题名称
	Group        string               // 消费者组
	Retry        int                  // 重试次数
	IsDelayQueue bool                 // 是否为延迟队列
	status       queue.ConsumerStatus // 消费者状态
	statusMu     sync.RWMutex         // 状态锁
	ctx          context.Context      // 上下文
	cancel       context.CancelFunc   // 取消函数
}

// Status 获取消费者状态
func (c *KafkaConsumer) Status() queue.ConsumerStatus {
	c.statusMu.RLock()
	defer c.statusMu.RUnlock()
	return c.status
}

// setStatus 设置消费者状态
func (c *KafkaConsumer) setStatus(status queue.ConsumerStatus) {
	c.statusMu.Lock()
	defer c.statusMu.Unlock()
	c.status = status
}

// Start 启动消费者
// h: 消息处理实现Handle(msg string) error方法
func (c *KafkaConsumer) Start(h queue.Handler) {
	c.setStatus(queue.ConsumerStatusRunning)
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go func() {
		defer c.setStatus(queue.ConsumerStatusStopped)

		for {
			select {
			case <-c.ctx.Done():
				c.Kafka.Log.Info(pkg.Sprintf("Kafka消费者 %s 已停止", c.Topic))
				return
			default:
				c.consumeLoop(h)
			}
		}
	}()
}

// consumeLoop 消费循环
func (c *KafkaConsumer) consumeLoop(h queue.Handler) {
	// 读取消息
	msg, err := c.Kafka.Reader.FetchMessage(c.ctx)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			c.Kafka.Log.Error(pkg.Sprintf("kafka fetch error: %v", err))
		}
		time.Sleep(time.Second)
		return
	}

	// 解析消息
	actualMsg := c.parseMessage(msg)

	// 重试处理
	var handleErr error
	for attempt := 0; attempt < c.Retry; attempt++ {
		handleErr = h.Handle(actualMsg)
		if handleErr == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if handleErr != nil {
		c.Kafka.Log.Error(pkg.Sprintf("kafka handle error after %d retries: %v", c.Retry, handleErr))
	}

	// 提交offset
	if err = c.Kafka.Reader.CommitMessages(c.ctx, msg); err != nil {
		c.Kafka.Log.Error(pkg.Sprintf("kafka commit error: %v", err))
	}
}

// parseMessage 解析消息
func (c *KafkaConsumer) parseMessage(msg kafka.Message) string {
	if c.IsDelayQueue {
		var msgMap map[string]any
		if err := json.Unmarshal(msg.Value, &msgMap); err != nil {
			c.Kafka.Log.Error(pkg.Sprintf("kafka delay msg unmarshal error: %v", err))
			return string(msg.Value)
		}

		// 检查是否到达消费时间
		if publishAt, ok := msgMap["publishAt"].(float64); ok {
			now := time.Now().UnixMilli()
			if now < int64(publishAt) {
				sleepMs := int64(publishAt) - now
				time.Sleep(time.Millisecond * time.Duration(sleepMs))
			}
		}

		if body, ok := msgMap["body"].(string); ok {
			return body
		}
		return string(msg.Value)
	}
	return string(msg.Value)
}

// Stop 停止消费者
func (c *KafkaConsumer) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}
	if c.Kafka.Reader != nil {
		return c.Kafka.Reader.Close()
	}
	return nil
}

// KafkaProducer Kafka生产者
// 封装了消息发布、延迟队列模拟、调试信息等通用逻辑
type KafkaProducer struct {
	Kafka        *Kafka // Kafka连接
	Topic        string // 主题名称
	Key          string // 消息Key
	IsDelayQueue bool   // 是否为延迟队列
	DelayMs      int64  // 延迟毫秒数
}

// Publish 发布消息到Kafka
// ctx: 上下文,用于链路追踪
// msg: 消息内容(字节数组)
func (p *KafkaProducer) Publish(ctx context.Context, msg []byte) error {
	start := time.Now()

	// 延迟队列模拟：包装消息添加发送时间戳
	if p.IsDelayQueue {
		msgMap := map[string]any{
			"body":      string(msg),
			"publishAt": time.Now().Add(time.Millisecond * time.Duration(p.DelayMs)).UnixMilli(),
		}
		msg, _ = json.Marshal(msgMap)
	}

	// 构建Kafka消息
	kmsg := kafka.Message{
		Value: msg,
	}
	if p.Key != "" {
		kmsg.Key = []byte(p.Key)
	}

	// 发送消息
	err := p.Kafka.Writer.WriteMessages(context.Background(), kmsg)
	if err != nil {
		p.Kafka.Log.Error(pkg.Sprintf("kafka publish error: %v", err))
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
	p.Kafka.Message.Publish(debugger.TopicMq, debugger.MqEvent{
		TraceId: traceId,
		Driver:  "kafka",
		Topic:   p.Topic,
		Message: string(msg),
		Ms:      float64(time.Since(start).Milliseconds()),
		Extra: map[string]interface{}{
			"key": p.Key,
			"err": err,
		},
	})

	return err
}

// Close 关闭生产者
func (p *KafkaProducer) Close() error {
	if p.Kafka.Writer != nil {
		return p.Kafka.Writer.Close()
	}
	return nil
}
