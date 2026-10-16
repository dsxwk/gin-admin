package queue

import (
	"context"
	"encoding/json"
	"errors"
	"gin/pkg/terminal"
	"gin/config"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/eventbus"
	"gin/pkg/serviceprovider/logger"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"fmt"
)

// Kafka Kafka连接
type Kafka struct {
	Writer *kafka.Writer
	Reader *kafka.Reader
	Conf   *config.Config
	Log    *logger.Logger
	Bus    *eventbus.Bus
}

// NewKafka 创建Kafka连接
func NewKafka(conf *config.Config, log *logger.Logger, bus *eventbus.Bus) *Kafka {
	return &Kafka{Conf: conf, Log: log, Bus: bus}
}

// KafkaConsumer Kafka消费者基类
type KafkaConsumer struct {
	Kafka    *Kafka
	Topic    string
	Group    string
	status   string
	statusMu sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// Status 获取消费者状态
func (c *KafkaConsumer) Status() string {
	c.statusMu.RLock()
	defer c.statusMu.RUnlock()
	if c.status == "" {
		return ConsumerStatusStopped
	}
	return c.status
}

// setStatus 设置消费者状态
func (c *KafkaConsumer) setStatus(status string) {
	c.statusMu.Lock()
	defer c.statusMu.Unlock()
	c.status = status
}

// Start 启动消费者
func (c *KafkaConsumer) Start(handler ConsumerHandler) {
	c.setStatus(ConsumerStatusStarting)
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go func() {
		defer c.setStatus(ConsumerStatusStopped)
		for {
			select {
			case <-c.ctx.Done():
				terminal.Infof("[Kafka] 消费者 %s 已停止", c.Topic)
				return
			default:
				c.consumeLoop(handler)
			}
		}
	}()
}

// consumeLoop 循环消费消息
func (c *KafkaConsumer) consumeLoop(handler ConsumerHandler) {
	if c.Kafka == nil || c.Kafka.Reader == nil {
		c.setStatus(ConsumerStatusError)
		time.Sleep(time.Second)
		return
	}
	c.setStatus(ConsumerStatusRunning)

	retry := handler.Retry()
	// FetchMessage 拉取消息
	msg, err := c.Kafka.Reader.FetchMessage(c.ctx)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			c.setStatus(ConsumerStatusError)
			c.Kafka.Log.Error(fmt.Sprintf("kafka fetch error: %v", err))
		}
		time.Sleep(time.Second)
		return
	}

	isDelay := handler.IsDelay()
	body := c.parseBody(msg, isDelay)

	if handleErr := handleWithRetry(c.ctx, handler, body, retry); handleErr != nil {
		c.Kafka.Log.Error(fmt.Sprintf("kafka handle error after %d retries: %v", retry, handleErr))
	}

	// CommitMessages 提交消息
	if err = c.Kafka.Reader.CommitMessages(c.ctx, msg); err != nil {
		c.setStatus(ConsumerStatusError)
		c.Kafka.Log.Error(fmt.Sprintf("kafka commit error: %v", err))
	}
}

// parseBody 解析消息内容
func (c *KafkaConsumer) parseBody(msg kafka.Message, isDelay bool) []byte {
	if isDelay {
		var msgMap map[string]any
		if err := json.Unmarshal(msg.Value, &msgMap); err != nil {
			c.Kafka.Log.Error(fmt.Sprintf("kafka delay msg unmarshal error: %v", err))
			return msg.Value
		}
		if body, ok := msgMap["body"].(string); ok {
			return []byte(body)
		}
		return msg.Value
	}
	return msg.Value
}

// Stop 停止消费者
func (c *KafkaConsumer) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}
	c.setStatus(ConsumerStatusStopped)
	if c.Kafka != nil && c.Kafka.Reader != nil {
		return c.Kafka.Reader.Close()
	}
	return nil
}

// KafkaProducer Kafka生产者基类
type KafkaProducer struct {
	Kafka *Kafka
	Topic string
	Key   string
	Owner Producer
}

// SetOwner 设置生产者归属
func (p *KafkaProducer) SetOwner(owner Producer) {
	p.Owner = owner
}

// Publish 发布消息
func (p *KafkaProducer) Publish(ctx context.Context, msg any) error {
	// 自动序列化为JSON
	body, err := marshalMessage(msg)
	if err != nil {
		return err
	}
	start := time.Now()

	if p.Owner != nil && p.Owner.IsDelay() {
		msgMap := map[string]any{
			"body":      string(body),
			"publishAt": time.Now().Add(time.Millisecond * time.Duration(p.Owner.DelayMs())).UnixMilli(),
		}
		body, _ = json.Marshal(msgMap)
	}

	kmsg := kafka.Message{Value: body}
	if p.Key != "" {
		kmsg.Key = []byte(p.Key)
	}

	// WriteMessages 写入消息
	err = p.Kafka.Writer.WriteMessages(context.Background(), kmsg)
	if err != nil {
		p.Kafka.Log.Error(fmt.Sprintf("kafka publish error: %v", err))
	}

	p.Kafka.Bus.Publish(ctx, debugger.TopicMQ, debugger.MQEvent{
		TraceID: traceID(ctx),
		Driver:  "kafka",
		Topic:   p.Topic,
		Message: string(body),
		Ms:      float64(time.Since(start).Milliseconds()),
		Extra: map[string]any{
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
