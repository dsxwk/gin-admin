package cache

import (
	"context"
	"fmt"
	"gin/common/ctxkey"
	"gin/config"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/eventbus"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
)

type redisHookContextKey int

const redisStartTimeKey redisHookContextKey = iota

type RedisHook struct {
	bus *eventbus.Bus
}

func (h *RedisHook) BeforeProcess(ctx context.Context, _ redis.Cmder) (context.Context, error) {
	// 在context中记录开始时间
	return context.WithValue(ctx, redisStartTimeKey, time.Now()), nil
}

func (h *RedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	start, ok := ctx.Value(redisStartTimeKey).(time.Time)
	if !ok {
		start = time.Now()
	}
	costMs := float64(time.Since(start).Nanoseconds()) / 1e6

	// 提取traceId
	traceId := "unknown"
	if id := ctx.Value(ctxkey.TraceIDKey); id != nil {
		if s, ok := id.(string); ok && s != "" {
			traceId = s
		}
	}

	// 发布事件
	if h.bus != nil {
		h.bus.Publish(ctx, debugger.TopicCache, debugger.CacheEvent{
			TraceID: traceId,
			Driver:  "redis",
			Name:    cmd.Name(),
			Cmd:     cmd.FullName(),
			Args:    cmd.Args(),
			Ms:      costMs,
		})
	}

	return nil
}

func (h *RedisHook) BeforeProcessPipeline(ctx context.Context, _ []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, redisStartTimeKey, time.Now()), nil
}

func (h *RedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	start, ok := ctx.Value(redisStartTimeKey).(time.Time)
	if !ok {
		start = time.Now()
	}
	costMs := float64(time.Since(start).Nanoseconds()) / 1e6

	// 提取traceId
	traceId := "unknown"
	if id := ctx.Value(ctxkey.TraceIDKey); id != nil {
		if s, ok := id.(string); ok && s != "" {
			traceId = s
		}
	}

	for _, cmd := range cmds {
		if h.bus != nil {
			h.bus.Publish(ctx, debugger.TopicCache, debugger.CacheEvent{
				TraceID: traceId,
				Driver:  "redis",
				Name:    cmd.Name(),
				Cmd:     cmd.FullName(),
				Args:    cmd.Args(),
				Ms:      costMs,
			})
		}
	}

	return nil
}

// redisState Redis共享状态
type redisState struct {
	mu        sync.RWMutex
	pubsubs   map[string]*redis.PubSub
	locks     map[string]*LockResult
	closed    bool
	closeOnce sync.Once
	closeErr  error
}

// RedisCache Redis缓存
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
	bus    *eventbus.Bus
	conf   *config.Config
	state  *redisState
}

var (
	redisCache   *CacheProxy
	redisCacheMu sync.Mutex
)

func NewRedisCache(conf *config.Config) *CacheProxy {
	redisCacheMu.Lock()
	defer redisCacheMu.Unlock()

	if redisCache != nil {
		return redisCache
	}

	bus := eventbus.Default()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Cache.Redis.Address, conf.Cache.Redis.Port),
		Password: conf.Cache.Redis.Password,
		DB:       conf.Cache.Redis.DB,
	})

	// 添加Hook
	client.AddHook(&RedisHook{bus: bus})

	r := &RedisCache{
		client: client,
		ctx:    context.Background(),
		bus:    bus,
		conf:   conf,
		state: &redisState{
			pubsubs: make(map[string]*redis.PubSub),
			locks:   make(map[string]*LockResult),
		},
	}

	redisCache = NewCacheProxy("redis", r, bus, conf)
	return redisCache
}

func (r *RedisCache) WithContext(ctx context.Context) *RedisCache {
	if r == nil {
		return nil
	}

	cp := *r
	cp.ctx = ctx

	return &cp
}

func (r *RedisCache) Set(key string, value any, expire time.Duration) error {
	data, err := encodeValue(value)
	if err != nil {
		return err
	}

	err = r.client.Set(r.ctx, key, data, expire).Err()
	if err != nil {
		return fmt.Errorf("error setting Redis cache: %w", err)
	}

	return nil
}

func (r *RedisCache) Get(key string) (any, bool) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, false
	}

	// 缓存只存储JSON,解析失败按未命中处理
	result, err := decodeValue([]byte(val))
	if err != nil {
		return nil, false
	}

	return result, true
}

func (r *RedisCache) Delete(key string) error {
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting Redis cache: %v", err)
	}

	return nil
}

// Exists 判断key是否存在
func (r *RedisCache) Exists(key string) (int64, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("error checking key exists: %v", err)
	}
	return result, nil
}

func (r *RedisCache) SAdd(key string, members ...any) error {
	values, err := encodeRedisSetMembers(members)
	if err != nil {
		return err
	}

	err = r.client.SAdd(r.ctx, key, values...).Err()
	if err != nil {
		return fmt.Errorf("error SAdd Redis set: %v", err)
	}
	return nil
}

func (r *RedisCache) SIsMember(key string, member any) (bool, error) {
	value, err := encodeRedisSetMember(member)
	if err != nil {
		return false, err
	}

	result, err := r.client.SIsMember(r.ctx, key, value).Result()
	if err != nil {
		return false, fmt.Errorf("error SIsMember Redis set: %v", err)
	}
	return result, nil
}

// encodeRedisSetMembers 编码Redis集合成员
func encodeRedisSetMembers(members []any) ([]any, error) {
	values := make([]any, len(members))
	for index, member := range members {
		value, err := encodeRedisSetMember(member)
		if err != nil {
			return nil, err
		}
		values[index] = value
	}

	return values, nil
}

// encodeRedisSetMember 编码Redis集合成员
func encodeRedisSetMember(member any) (any, error) {
	switch member.(type) {
	case string, []byte:
		return member, nil
	default:
		value, err := json.Marshal(member)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal set member: %w", err)
		}

		return string(value), nil
	}
}

// GetWithTTL 获取缓存值和过期时间
func (r *RedisCache) GetWithTTL(key string) (any, time.Time, bool, error) {
	ttl, err := r.client.PTTL(r.ctx, key).Result()
	if err != nil {
		return nil, time.Time{}, false, fmt.Errorf("error getting TTL for key %v: %w", key, err)
	}

	if ttl == -2*time.Nanosecond {
		return nil, time.Time{}, false, nil
	}

	value, ok := r.Get(key)
	if !ok {
		return nil, time.Time{}, false, nil
	}

	if ttl == -1*time.Nanosecond {
		return value, time.Time{}, true, nil
	}

	return value, time.Now().Add(ttl), true, nil
}

// Expire 获取缓存值和过期时间
func (r *RedisCache) Expire(key string) (any, time.Time, bool, error) {
	return r.GetWithTTL(key)
}

// Publish 发布
func (r *RedisCache) Publish(channel string, message any) error {
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

	r.state.mu.Lock()
	if r.state.closed {
		r.state.mu.Unlock()
		_ = pubsub.Close()
		return ErrCacheClosed
	}
	if _, exists := r.state.pubsubs[channel]; exists {
		r.state.mu.Unlock()
		_ = pubsub.Close()
		return fmt.Errorf("channel %s already subscribed", channel)
	}
	r.state.pubsubs[channel] = pubsub
	r.state.mu.Unlock()

	// 消息处理协程
	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			handler(msg.Channel, msg.Payload)
		}
	}()

	return nil
}

// Pipeline 返回Redis Pipeline,仅Redis驱动可用
func (r *RedisCache) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

// Ping 检查Redis连接是否正常
func (r *RedisCache) Ping() error {
	return r.client.Ping(r.ctx).Err()
}

// Client 获取Redis客户端
func (r *RedisCache) Client() *redis.Client {
	return r.client
}

// Unsubscribe 取消订阅
func (r *RedisCache) Unsubscribe(channel string) error {
	r.state.mu.Lock()
	pubsub, ok := r.state.pubsubs[channel]
	if !ok {
		r.state.mu.Unlock()
		return fmt.Errorf("channel %s not found in subscriptions", channel)
	}
	delete(r.state.pubsubs, channel)
	r.state.mu.Unlock()

	err := pubsub.Unsubscribe(r.ctx, channel)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe from channel %s: %v", channel, err)
	}

	// 关闭并删除记录
	err = pubsub.Close()
	if err != nil {
		return fmt.Errorf("failed to close pubsub for channel %s: %v", channel, err)
	}

	return nil
}

// Close 关闭Redis客户端
func (r *RedisCache) Close() error {
	if r == nil || r.state == nil || r.client == nil {
		return nil
	}

	r.state.closeOnce.Do(func() {
		r.state.mu.Lock()
		r.state.closed = true
		pubsubs := r.state.pubsubs
		locks := r.state.locks
		r.state.pubsubs = make(map[string]*redis.PubSub)
		r.state.locks = make(map[string]*LockResult)
		r.state.mu.Unlock()

		for _, lock := range locks {
			_ = lock.Release()
		}
		for _, pubsub := range pubsubs {
			_ = pubsub.Close()
		}

		r.state.closeErr = r.client.Close()
	})

	return r.state.closeErr
}
