package base

import (
	"context"
	"fmt"
	"gin/app/facade"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/serviceprovider/queue"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"sync"
	"time"
)

type RedisConsumer struct {
	Queue    string
	status   queue.ConsumerStatus
	statusMu sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

func (c *RedisConsumer) Status() queue.ConsumerStatus {
	c.statusMu.RLock()
	defer c.statusMu.RUnlock()
	return c.status
}

func (c *RedisConsumer) setStatus(status queue.ConsumerStatus) {
	c.statusMu.Lock()
	defer c.statusMu.Unlock()
	c.status = status
}

func (c *RedisConsumer) Start(h interface{}) {
	c.setStatus(queue.ConsumerStatusRunning)
	c.ctx, c.cancel = context.WithCancel(context.Background())
	go func() {
		defer c.setStatus(queue.ConsumerStatusStopped)
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

func (c *RedisConsumer) consumeLoop(h interface{}) {
	client := facade.Cache("redis").Redis().Client()
	if client == nil {
		time.Sleep(time.Second)
		return
	}
	if h.(queue.Consumer).IsDelay() {
		c.processDelayed(client, h)
		time.Sleep(time.Second)
		return
	}
	c.processNormal(client, h)
}

func (c *RedisConsumer) processNormal(client *redis.Client, h interface{}) {
	result, err := client.BRPop(c.ctx, 3*time.Second, c.Queue).Result()
	if err != nil {
		return
	}
	c.handleMessage([]byte(result[1]), h)
}

func (c *RedisConsumer) processDelayed(client *redis.Client, h interface{}) {
	delayedKey := c.Queue + ":delayed"
	now := float64(time.Now().UnixMilli())
	members, err := client.ZRangeByScoreWithScores(c.ctx, delayedKey, &redis.ZRangeBy{
		Min: "0", Max: pkg.Sprintf("%.0f", now), Count: 10,
	}).Result()
	if err != nil || len(members) == 0 {
		return
	}
	for _, m := range members {
		c.handleMessage([]byte(m.Member.(string)), h)
		client.ZRem(c.ctx, delayedKey, m.Member)
	}
}

func (c *RedisConsumer) handleMessage(body []byte, h interface{}) {
	retry := h.(queue.Consumer).Retry()
	var handleErr error
	for attempt := 0; attempt < retry || attempt == 0; attempt++ {
		handleErr = queue.TryHandle(h, body)
		if handleErr == nil {
			return
		}
		if retry > 0 && attempt < retry {
			time.Sleep(time.Second)
		}
	}
	if handleErr != nil {
		facade.Log().Error(pkg.Sprintf("[Redis] handle error: %v", handleErr))
	}
}

func (c *RedisConsumer) Stop() error {
	if c.cancel != nil {
		c.cancel()
	}
	return nil
}

type RedisProducer struct {
	Queue string
	Owner queue.Producer
}

func (p *RedisProducer) Publish(ctx context.Context, msg any) error {
	client := facade.Cache("redis").Redis().Client()
	if client == nil {
		return fmt.Errorf("redis client not initialized")
	}
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
	if p.Owner != nil && p.Owner.IsDelay() && p.Owner.DelayMs() > 0 {
		delayedKey := p.Queue + ":delayed"
		score := float64(time.Now().UnixMilli() + p.Owner.DelayMs())
		return client.ZAdd(ctx, delayedKey, &redis.Z{Score: score, Member: string(body)}).Err()
	}
	return client.LPush(ctx, p.Queue, string(body)).Err()
}

func (p *RedisProducer) Close() error { return nil }
