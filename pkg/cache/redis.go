package cache

import (
	"context"
	"fmt"
	"gin/config"
	"gin/pkg/debugger"
	"gin/pkg/message"
	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"sync"
	"time"
)

type RedisHook struct {
	bus *message.EventBus
}

func (h *RedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	// 在context中记录开始时间
	return context.WithValue(ctx, "startTime", time.Now()), nil
}

func (h *RedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	start, ok := ctx.Value("startTime").(time.Time)
	if !ok {
		start = time.Now()
	}
	costMs := float64(time.Since(start).Nanoseconds()) / 1e6

	// 发布事件
	if h.bus != nil {
		h.bus.Publish(debugger.TopicCache, debugger.CacheEvent{
			Driver: "redis",
			Name:   cmd.Name(),
			Cmd:    cmd.FullName(),
			Args:   cmd.Args(),
			Ms:     costMs,
		})
	}

	return nil
}

func (h *RedisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, "startTime", time.Now()), nil
}

func (h *RedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	start, ok := ctx.Value("startTime").(time.Time)
	if !ok {
		start = time.Now()
	}
	costMs := float64(time.Since(start).Nanoseconds()) / 1e6

	for _, cmd := range cmds {
		if h.bus != nil {
			h.bus.Publish(debugger.TopicCache, debugger.CacheEvent{
				Driver: "redis",
				Name:   cmd.Name(),
				Cmd:    cmd.FullName(),
				Args:   cmd.Args(),
				Ms:     costMs,
			})
		}
	}

	return nil
}

// RedisCache Redis缓存
type RedisCache struct {
	client  *redis.Client
	pubsubs map[string]*redis.PubSub
	ctx     context.Context
	bus     *message.EventBus
}

var (
	redisCache *CacheProxy
	redisOnce  sync.Once
)

func NewRedisCache(conf *config.Config) *CacheProxy {
	redisOnce.Do(func() {
		var (
			bus = message.GetEventBus()
		)

		client := redis.NewClient(&redis.Options{
			Addr:     conf.Cache.Redis.Address,
			Password: conf.Cache.Redis.Password,
			DB:       conf.Cache.Redis.DB,
		})

		// 添加Hook
		client.AddHook(&RedisHook{bus: bus})

		r := &RedisCache{
			client:  client,
			ctx:     context.Background(),
			pubsubs: make(map[string]*redis.PubSub),
			bus:     bus,
		}

		redisCache = NewCacheProxy("redis", r, bus)
	})
	return redisCache
}

func (r *RedisCache) WithContext(ctx context.Context) *RedisCache {
	return &RedisCache{
		client:  r.client,
		pubsubs: r.pubsubs,
		ctx:     ctx,
	}
}

func (r *RedisCache) Set(key string, value interface{}, expire time.Duration) error {
	var valStr string

	// 根据类型处理值
	switch v := value.(type) {
	case string:
		valStr = v
	case []byte:
		valStr = string(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		valStr = fmt.Sprintf("%v", v)
	case float32, float64:
		valStr = fmt.Sprintf("%v", v)
	case bool:
		valStr = fmt.Sprintf("%v", v)
	default:
		// map,slice,struct等复杂类型需要JSON序列化
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %v", err)
		}
		valStr = string(b)
	}

	err := r.client.Set(r.ctx, key, valStr, expire).Err()
	if err != nil {
		return fmt.Errorf("error setting Redis cache: %v", err)
	}

	return nil
}

func (r *RedisCache) Get(key string) (interface{}, bool) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, false
	}

	// 尝试解析JSON
	var result interface{}
	if err = json.Unmarshal([]byte(val), &result); err == nil {
		// 如果是JSON对象或数组,返回解析后的结果
		return result, true
	}

	// 否则返回原始字符串
	return val, true
}

func (r *RedisCache) Delete(key string) error {
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting Redis cache: %v", err)
	}

	return nil
}

func (r *RedisCache) Expire(key string) (interface{}, time.Time, bool, error) {
	// Redis不支持使用相同的API获取到期时间,因此必须使用TTL
	ttl, err := r.client.TTL(r.ctx, key).Result()
	if err != nil {
		return nil, time.Time{}, false, fmt.Errorf("error getting TTL for key %v: %v", key, err)
	}

	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, time.Time{}, false, fmt.Errorf("error getting value for key %v: %v", key, err)
	}

	expireTime := time.Now().Add(ttl)

	return val, expireTime, true, nil
}

// Lock 获取锁
func (r *RedisCache) Lock(key string, value string, expire time.Duration) error {
	// 使用 SETNX 命令尝试设置锁
	result, err := r.client.SetNX(r.ctx, key, value, expire).Result()
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %v", err)
	}

	if !result {
		// 如果返回 false,表示锁已存在
		return fmt.Errorf("lock already exists")
	}

	return nil
}

// UnLock 释放锁
func (r *RedisCache) UnLock(key string, value string) error {
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end`
	// 使用 EVAL 命令执行 Lua 脚本
	status, err := r.client.Eval(r.ctx, script, []string{key}, value).Int()
	if err != nil {
		return fmt.Errorf("failed to unlock: %v", err)
	}

	if status == 0 {
		return fmt.Errorf("unlock failed: lock not owned or already released")
	}

	return nil
}

// Publish 发布
func (r *RedisCache) Publish(channel string, message interface{}) error {
	var (
		payload string
	)

	switch v := message.(type) {
	case string:
		payload = v
	case []byte:
		payload = string(v)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %v", err)
		}
		payload = string(b)
	}

	err := r.client.Publish(r.ctx, channel, payload).Err()
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}

// Subscribe 订阅
func (r *RedisCache) Subscribe(channel string, handler func(channel string, payload string)) error {
	pubsub := r.client.Subscribe(r.ctx, channel)

	// 等待订阅确认
	_, err := pubsub.Receive(r.ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe to channel %s: %v", channel, err)
	}

	// 保存pubsub对象
	r.pubsubs[channel] = pubsub

	// 消息处理协程
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			handler(msg.Channel, msg.Payload)
		}
	}()

	return nil
}

// Unsubscribe 取消订阅
func (r *RedisCache) Unsubscribe(channel string) error {
	pubsub, ok := r.pubsubs[channel]
	if !ok {
		return fmt.Errorf("channel %s not found in subscriptions", channel)
	}

	err := pubsub.Unsubscribe(r.ctx, channel)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe from channel %s: %v", channel, err)
	}

	// 关闭并删除记录
	err = pubsub.Close()
	if err != nil {
		return fmt.Errorf("failed to close pubsub for channel %s: %v", channel, err)
	}

	delete(r.pubsubs, channel)
	return nil
}
