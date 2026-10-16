package queue

import (
	"context"
	"encoding/json"
	"gin/pkg/ctxkey"
	"time"
)

// marshalMessage 统一序列化生产者消息
func marshalMessage(msg any) ([]byte, error) {
	switch value := msg.(type) {
	case []byte:
		return value, nil
	case string:
		return []byte(value), nil
	default:
		return json.Marshal(msg)
	}
}

// traceID 获取请求链路标识
func traceID(ctx context.Context) string {
	if ctx == nil {
		return "unknown"
	}
	if id, ok := ctx.Value(ctxkey.TraceIDKey).(string); ok && id != "" {
		return id
	}
	return "unknown"
}

// handleWithRetry 自动反序列化并按指定次数重试
func handleWithRetry(ctx context.Context, handler ConsumerHandler, body []byte, retries int) error {
	if retries < 1 {
		retries = 1
	}

	var handleErr error
	for attempt := 0; attempt < retries; attempt++ {
		if handleErr = TryHandle(ctx, handler, body); handleErr == nil {
			return nil
		}
		if attempt+1 < retries {
			time.Sleep(time.Second)
		}
	}
	return handleErr
}

// TryHandle 自动反序列化并调用处理,支持上下文
func TryHandle(ctx context.Context, handler PayloadHandler, body []byte) error {
	if ctx == nil {
		ctx = context.Background()
	}
	payload := handler.NewPayload()
	if err := json.Unmarshal(body, payload); err != nil {
		return err
	}
	if contextual, ok := handler.(ContextPayloadHandler); ok {
		return contextual.HandleContext(ctx, payload)
	}
	return handler.Handle(payload)
}
