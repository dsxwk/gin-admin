package job

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Execute 执行任务消息
func Execute(ctx context.Context, item Job, message Message) error {
	if item == nil {
		return fmt.Errorf("job [%s] 未注册", message.JobName)
	}
	if err := waitUntil(ctx, message.RunAt); err != nil {
		return err
	}

	payload := item.NewPayload()
	if err := json.Unmarshal(message.Payload, payload); err != nil {
		return fmt.Errorf("job [%s] payload解析失败: %w", message.JobName, err)
	}

	attempts := item.Retry() + 1
	if attempts < 1 {
		attempts = 1
	}
	retryDelay := time.Duration(item.Delay()) * time.Millisecond

	var lastErr error
	for attempt := 0; attempt < attempts; attempt++ {
		if lastErr = item.Handle(payload); lastErr == nil {
			return nil
		}
		if attempt+1 < attempts {
			if err := sleepContext(ctx, retryDelay); err != nil {
				return err
			}
		}
	}
	return fmt.Errorf("job [%s] 处理失败: %w", message.JobName, lastErr)
}

// waitUntil 等待任务执行时间
func waitUntil(ctx context.Context, runAt int64) error {
	if runAt <= 0 {
		return nil
	}
	delay := time.Until(time.UnixMilli(runAt))
	return sleepContext(ctx, delay)
}

// sleepContext 支持上下文取消的等待
func sleepContext(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
