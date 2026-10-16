package queue

import (
	"context"
	"errors"
	"fmt"
	"gin/pkg/terminal"
	"gin/pkg/serviceprovider/logger"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisConsumer Redis消费者基类
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

// Status 获取消费者状态
func (c *RedisConsumer) Status() string {
	c.statusMu.RLock()
	defer c.statusMu.RUnlock()
	if c.status == "" {
		return ConsumerStatusStopped
	}
	return c.status
}

// setStatus 设置消费者状态
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

// Start 启动消费者
func (c *RedisConsumer) Start(handler ConsumerHandler) {
	c.setStatus(ConsumerStatusStarting)
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go func() {
		defer c.setStatus(ConsumerStatusStopped)
		for {
			select {
			case <-c.ctx.Done():
				terminal.Infof("[Redis] consumer %s stopped", c.Queue)
				return
			default:
				c.consumeLoop(handler)
			}
		}
	}()
}

// consumeLoop 循环消费消息
func (c *RedisConsumer) consumeLoop(handler ConsumerHandler) {
	client := c.client()
	if client == nil {
		c.setStatus(ConsumerStatusError)
		time.Sleep(time.Second)
		return
	}
	c.setStatus(ConsumerStatusRunning)
	if handler.IsDelay() || c.Dual {
		c.processDelayed(client, handler)
	}
	if handler.IsDelay() && !c.Dual {
		time.Sleep(time.Second)
		return
	}
	if err := c.processNormal(client, handler); err != nil {
		time.Sleep(time.Second)
	}
}

// processNormal 处理普通队列消息
func (c *RedisConsumer) processNormal(client *redis.Client, handler ConsumerHandler) error {
	result, err := client.BRPop(c.ctx, 3*time.Second, c.Queue).Result()
	if err != nil {
		if c.ctx.Err() != nil || errors.Is(err, redis.Nil) {
			return nil
		}
		c.setStatus(ConsumerStatusError)
		return err
	}
	c.setStatus(ConsumerStatusRunning)
	c.handleMessage([]byte(result[1]), handler)
	return nil
}

// processDelayed 处理延迟队列消息
func (c *RedisConsumer) processDelayed(client *redis.Client, handler ConsumerHandler) {
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
		c.handleMessage(DecodeRedisDelayedMessage(member), handler)
	}
}

// handleMessage 处理消息
func (c *RedisConsumer) handleMessage(body []byte, handler ConsumerHandler) {
	if handleErr := handleWithRetry(c.ctx, handler, body, handler.Retry()); handleErr != nil {
		if c.Log != nil {
			c.Log.Error(fmt.Sprintf("[Redis] handle error: %v", handleErr))
		}
	}
}

// Stop 停止消费者
func (c *RedisConsumer) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}
	c.setStatus(ConsumerStatusStopped)
	return nil
}

// RedisProducer Redis生产者基类
type RedisProducer struct {
	Queue        string
	DelayedQueue string
	GetClient    func() *redis.Client
	Owner        Producer
}

// SetOwner 设置生产者归属
func (p *RedisProducer) SetOwner(owner Producer) {
	p.Owner = owner
}

// Publish 发布消息
func (p *RedisProducer) Publish(ctx context.Context, msg any) error {
	body, err := marshalMessage(msg)
	if err != nil {
		return err
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

// Close 关闭生产者
func (p *RedisProducer) Close() error { return nil }
