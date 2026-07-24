package base

import (
	"context"
	"errors"
	"gin/common/ctxkey"
	"gin/common/flag"
	"gin/config"
	"gin/pkg"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/logger"
	"gin/pkg/serviceprovider/message"
	"gin/pkg/serviceprovider/queue"
	"github.com/goccy/go-json"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

// Kafka Kafka连接
type Kafka struct {
	Writer  *kafka.Writer
	Reader  *kafka.Reader
	Conf    *config.Config
	Log     *logger.Logger
	Message *message.Event
}

func NewKafka(conf *config.Config, log *logger.Logger, bus *message.Event) *Kafka {
	return &Kafka{Conf: conf, Log: log, Message: bus}
}

// KafkaConsumer Kafka消费者基类
type KafkaConsumer struct {
	Kafka    *Kafka
	Topic    string
	Group    string
	status   queue.ConsumerStatus
	statusMu sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

func (c *KafkaConsumer) Status() queue.ConsumerStatus {
	c.statusMu.RLock()
	defer c.statusMu.RUnlock()
	return c.status
}

func (c *KafkaConsumer) setStatus(status queue.ConsumerStatus) {
	c.statusMu.Lock()
	defer c.statusMu.Unlock()
	c.status = status
}

func (c *KafkaConsumer) Start(h interface{}) {
	c.setStatus(queue.ConsumerStatusRunning)
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go func() {
		defer c.setStatus(queue.ConsumerStatusStopped)
		for {
			select {
			case <-c.ctx.Done():
				flag.Infof("[Kafka] 消费者 %s 已停止", c.Topic)
				return
			default:
				c.consumeLoop(h)
			}
		}
	}()
}

func (c *KafkaConsumer) consumeLoop(h interface{}) {
	retry := h.(queue.Consumer).Retry()
	msg, err := c.Kafka.Reader.FetchMessage(c.ctx)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			c.Kafka.Log.Error(pkg.Sprintf("kafka fetch error: %v", err))
		}
		time.Sleep(time.Second)
		return
	}

	isDelay := h.(queue.Consumer).IsDelay()
	body := c.parseBody(msg, isDelay)

	var handleErr error
	for attempt := 0; attempt < retry; attempt++ {
		handleErr = queue.TryHandle(h, body)
		if handleErr == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if handleErr != nil {
		c.Kafka.Log.Error(pkg.Sprintf("kafka handle error after %d retries: %v", retry, handleErr))
	}

	if err = c.Kafka.Reader.CommitMessages(c.ctx, msg); err != nil {
		c.Kafka.Log.Error(pkg.Sprintf("kafka commit error: %v", err))
	}
}

func (c *KafkaConsumer) parseBody(msg kafka.Message, isDelay bool) []byte {
	if isDelay {
		var msgMap map[string]any
		if err := json.Unmarshal(msg.Value, &msgMap); err != nil {
			c.Kafka.Log.Error(pkg.Sprintf("kafka delay msg unmarshal error: %v", err))
			return msg.Value
		}
		if body, ok := msgMap["body"].(string); ok {
			return []byte(body)
		}
		return msg.Value
	}
	return msg.Value
}

func (c *KafkaConsumer) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}
	if c.Kafka.Reader != nil {
		return c.Kafka.Reader.Close()
	}
	return nil
}

// KafkaProducer Kafka生产者基类
type KafkaProducer struct {
	Kafka *Kafka
	Topic string
	Key   string
	Owner queue.Producer
}

func (p *KafkaProducer) Publish(ctx context.Context, msg any) error {
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

	err := p.Kafka.Writer.WriteMessages(context.Background(), kmsg)
	if err != nil {
		p.Kafka.Log.Error(pkg.Sprintf("kafka publish error: %v", err))
	}

	traceId := "unknown"
	if ctx != nil {
		if id := ctx.Value(ctxkey.TraceIdKey); id != nil {
			if s, ok := id.(string); ok && s != "" {
				traceId = s
			}
		}
	}

	p.Kafka.Message.Publish(debugger.TopicMq, debugger.MqEvent{
		TraceId: traceId,
		Driver:  "kafka",
		Topic:   p.Topic,
		Message: string(body),
		Ms:      float64(time.Since(start).Milliseconds()),
		Extra: map[string]interface{}{
			"key": p.Key,
			"err": err,
		},
	})

	return err
}

func (p *KafkaProducer) Close() error {
	if p.Kafka.Writer != nil {
		return p.Kafka.Writer.Close()
	}
	return nil
}
