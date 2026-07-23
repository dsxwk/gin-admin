package work

import (
	"context"
	"encoding/json"
	"fmt"
	"gin/app/facade"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/serviceprovider/job"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	JobQueueKey   = "job:queue"         // Redis任务队列key
	JobDelayedKey = "job:queue:delayed" // Redis延迟队列key
	JobPopTimeout = 3 * time.Second     // BRPop阻塞超时
)

// RedisWorker Redis任务消费者
type RedisWorker struct {
	ctx    context.Context
	cancel context.CancelFunc
}

// NewRedisWorker 创建RedisWorker
func NewRedisWorker() *RedisWorker {
	return &RedisWorker{}
}

// Start 启动消费者
func (w *RedisWorker) Start() error {
	w.ctx, w.cancel = context.WithCancel(context.Background())

	go func() {
		flag.Infof("Job Redis Worker 已启动")
		for {
			select {
			case <-w.ctx.Done():
				flag.Infof("Job Redis Worker 已停止")
				return
			default:
				w.consume()
			}
		}
	}()
	return nil
}

// Stop 停止消费者
func (w *RedisWorker) Stop() error {
	if w.cancel != nil {
		w.cancel()
	}
	return nil
}

// consume 消费循环
func (w *RedisWorker) consume() {
	client := facade.Cache("redis").Redis().Client()
	if client == nil {
		time.Sleep(time.Second)
		return
	}

	// 优先检查延迟队列
	w.processDelayed(client)

	// 再检查普通队列
	result, err := client.BRPop(w.ctx, JobPopTimeout, JobQueueKey).Result()
	if err != nil {
		return
	}

	var msg JobMessage
	if err = json.Unmarshal([]byte(result[1]), &msg); err != nil {
		facade.Log().Error(pkg.Sprintf("Job 消息解析失败: %v", err))
		return
	}

	w.handleMessage(msg)
}

// processDelayed 处理延迟队列中的就绪任务
func (w *RedisWorker) processDelayed(client *redis.Client) {
	now := float64(time.Now().UnixMilli())

	members, err := client.ZRangeByScoreWithScores(w.ctx, JobDelayedKey, &redis.ZRangeBy{
		Min:   "0",
		Max:   fmt.Sprintf("%.0f", now),
		Count: 10,
	}).Result()
	if err != nil || len(members) == 0 {
		return
	}

	for _, m := range members {
		var msg JobMessage
		if err = json.Unmarshal([]byte(m.Member.(string)), &msg); err != nil {
			continue
		}
		w.handleMessage(msg)
		client.ZRem(w.ctx, JobDelayedKey, m.Member)
	}
}

// handleMessage 处理消息
func (w *RedisWorker) handleMessage(msg JobMessage) {
	j := job.Get(msg.JobName)
	if j == nil {
		facade.Log().Error(pkg.Sprintf("Job [%s] 未注册", msg.JobName))
		return
	}

	// 延迟执行
	if msg.RunAt > 0 {
		now := time.Now().UnixMilli()
		if msg.RunAt > now {
			time.Sleep(time.Millisecond * time.Duration(msg.RunAt-now))
		}
	}

	payload := j.NewPayload()
	if err := json.Unmarshal(msg.Payload, payload); err != nil {
		facade.Log().Error(pkg.Sprintf("Job [%s] payload解析失败: %v", msg.JobName, err))
		return
	}

	retry := j.Retry()
	delay := j.Delay()

	var err error
	for attempt := 0; attempt < retry || attempt == 0; attempt++ {
		err = j.Handle(payload)
		if err == nil {
			return
		}
		facade.Log().Error(pkg.Sprintf("Job [%s] 处理失败(attempt %d/%d): %v", msg.JobName, attempt+1, retry, err))
		if delay > 0 {
			time.Sleep(time.Millisecond * time.Duration(delay))
		}
	}
}
