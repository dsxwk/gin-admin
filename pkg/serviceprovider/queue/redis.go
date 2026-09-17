package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisMessageID atomic.Uint64

// RedisDelayedMessage Redis延迟消息信封
type RedisDelayedMessage struct {
	ID   string          `json:"id"`
	Body json.RawMessage `json:"body"`
}

// NewRedisDelayedMessage 创建带唯一标识的延迟消息
func NewRedisDelayedMessage(body []byte) ([]byte, error) {
	message := RedisDelayedMessage{
		ID:   fmt.Sprintf("%d-%d", time.Now().UnixNano(), redisMessageID.Add(1)),
		Body: body,
	}
	return json.Marshal(message)
}

// DecodeRedisDelayedMessage 解析延迟消息,兼容旧格式
func DecodeRedisDelayedMessage(member string) []byte {
	var message RedisDelayedMessage
	if err := json.Unmarshal([]byte(member), &message); err == nil && message.ID != "" && len(message.Body) > 0 {
		return message.Body
	}
	return []byte(member)
}

const claimRedisDelayedScript = `
local messages = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1], 'LIMIT', 0, ARGV[2])
if #messages == 0 then
	return {}
end
redis.call('ZREM', KEYS[1], unpack(messages))
return messages
`

var claimRedisDelayed = redis.NewScript(claimRedisDelayedScript)

// ClaimRedisDelayedMessages 原子领取到期延迟消息
func ClaimRedisDelayedMessages(
	ctx context.Context,
	client *redis.Client,
	key string,
	maxScore float64,
	count int64,
) ([]string, error) {
	return claimRedisDelayed.Run(
		ctx,
		client,
		[]string{key},
		fmt.Sprintf("%.0f", maxScore),
		count,
	).StringSlice()
}

// PublishRedisMessage 发布Redis普通或延迟消息
func PublishRedisMessage(
	ctx context.Context,
	client *redis.Client,
	queue string,
	delayedQueue string,
	body []byte,
	delay time.Duration,
) error {
	if client == nil {
		return fmt.Errorf("redis客户端未初始化")
	}
	if delay <= 0 {
		return client.LPush(ctx, queue, string(body)).Err()
	}

	delayedBody, err := NewRedisDelayedMessage(body)
	if err != nil {
		return err
	}
	score := float64(time.Now().Add(delay).UnixMilli())
	return client.ZAdd(ctx, delayedQueue, &redis.Z{
		Score:  score,
		Member: string(delayedBody),
	}).Err()
}
