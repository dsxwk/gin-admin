package queue

import (
	"context"
	"errors"
	"fmt"
	"gin/common/flag"
	"gin/pkg/serviceprovider/logger"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
)

type RedisConsumer struct {
	Queue        string
	DelayedQueue string
	GetClient    func() *redis.Client
	Log          *logger.Logger
	Dual         bool
	status       string
	statusMu     sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
}

func (c *RedisConsumer) Status() string {
	c.statusMu.RLock()
	defer c.statusMu.RUnlock()
	if c.status == "" {
		return ConsumerStatusStopped
	}
	return c.status
}

func (c *RedisConsumer) setStatus(status string) {
	c.statusMu.Lock()
	defer c.statusMu.Unlock()
	c.status = status
}

// client 获取Redis客户端
func (c *RedisConsumer) client() *redis.Client {
	if c.GetClient == nil {
		return nil
	}
	return c.GetClient()
}

// delayedQueue 获取延迟队列名称
func (c *RedisConsumer) delayedQueue() string {
	if c.DelayedQueue != "" {
		return c.DelayedQueue
	}
	return c.Queue + ":delayed"
}

func (c *RedisConsumer) Start[T ConsumerHandler](h T) {
	c.setStatus(ConsumerStatusStarting)
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go func() {
		defer c.setStatus(ConsumerStatusStopped)
		for {
			select {
			case <-c.ctx.Done():
				flag.Infof("[Redis] consumer %s stopped", c.Queue)
				return
			default:
				c.consumeLoop(h)
			}
		}
	}()
}

func (c *RedisConsumer) consumeLoop[T ConsumerHandler](h T) {
	client := c.client()
	if client == nil {
		c.setStatus(ConsumerStatusError)
		time.Sleep(time.Second)
		return
	}
	c.setStatus(ConsumerStatusRunning)
	if h.IsDelay() || c.Dual {
		c.processDelayed(client, h)
	}
	if h.IsDelay() && !c.Dual {
		time.Sleep(time.Second)
		return
	}
	if err := c.processNormal(client, h); err != nil {
		time.Sleep(time.Second)
	}
}

func (c *RedisConsumer) processNormal[T ConsumerHandler](client *redis.Client, h T) error {
	result, err := client.BRPop(c.ctx, 3*time.Second, c.Queue).Result()
	if err != nil {
		if c.ctx.Err() != nil || errors.Is(err, redis.Nil) {
			return nil
		}
		c.setStatus(ConsumerStatusError)
		return err
	}
	c.setStatus(ConsumerStatusRunning)
	c.handleMessage([]byte(result[1]), h)
	return nil
}

func (c *RedisConsumer) processDelayed[T ConsumerHandler](client *redis.Client, h T) {
	delayedKey := c.delayedQueue()
	now := float64(time.Now().UnixMilli())
	members, err := ClaimRedisDelayedMessages(c.ctx, client, delayedKey, now, 10)
	if err != nil {
		if c.ctx.Err() == nil {
			c.setStatus(ConsumerStatusError)
		}
		return
	}
	if len(members) == 0 {
		return
	}
	c.setStatus(ConsumerStatusRunning)
	for _, member := range members {
		c.handleMessage(DecodeRedisDelayedMessage(member), h)
	}
}

func (c *RedisConsumer) handleMessage[T ConsumerHandler](body []byte, h T) {
	retry := h.Retry()
	var handleErr error
	for attempt := 0; attempt < retry || attempt == 0; attempt++ {
		handleErr = TryHandle(c.ctx, h, body)
		if handleErr == nil {
			return
		}
		if retry > 0 && attempt < retry {
			time.Sleep(time.Second)
		}
	}
	if handleErr != nil {
		if c.Log != nil {
			c.Log.Error(fmt.Sprintf("[Redis] handle error: %v", handleErr))
		}
	}
}

func (c *RedisConsumer) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}
	c.setStatus(ConsumerStatusStopped)
	return nil
}

type RedisProducer struct {
	Queue        string
	DelayedQueue string
	GetClient    func() *redis.Client
	Owner        Producer
}

func (p *RedisProducer) Publish(ctx context.Context, msg any) error {
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
	var delay time.Duration
	if p.Owner != nil && p.Owner.IsDelay() && p.Owner.DelayMs() > 0 {
		delay = time.Duration(p.Owner.DelayMs()) * time.Millisecond
	}
	return p.PublishRaw(ctx, body, delay)
}

// PublishRaw 发布原始Redis消息
func (p *RedisProducer) PublishRaw(ctx context.Context, body []byte, delay time.Duration) error {
	if p.GetClient == nil {
		return fmt.Errorf("redis client not initialized")
	}
	client := p.GetClient()
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	delayedQueue := p.DelayedQueue
	if delayedQueue == "" {
		delayedQueue = p.Queue + ":delayed"
	}
	return PublishRedisMessage(ctx, client, p.Queue, delayedQueue, body, delay)
}

func (p *RedisProducer) Close() error { return nil }
