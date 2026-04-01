package tests

import (
	"context"
	"gin/app/facade"
	"gin/common/ctxkey"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestRabbitMQPublish 测试RabbitMQ普通消息发布和消费
func TestRabbitMQPublish(t *testing.T) {
	cfg := facade.Config.Get()
	if !cfg.Rabbitmq.Enabled {
		t.Skip("RabbitMQ not enabled, skipping test")
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxkey.TraceIdKey, "test-trace-id")

	// 等待消费者启动
	time.Sleep(500 * time.Millisecond)

	// 获取生产者
	producer := facade.Queue.Producer("rabbitmq_demo")
	if producer == nil {
		t.Skip("RabbitMQ producer not registered")
	}
	defer func() {
		err := producer.Close()
		if err != nil {
			t.Errorf("Failed to close producer: %v", err)
		}
	}()

	// 发送测试消息
	testMessages := []string{
		`{"orderId":1, "amount":100, "product":"apple"}`,
		`{"orderId":2, "amount":200, "product":"banana"}`,
		`{"orderId":3, "amount":300, "product":"orange"}`,
	}

	for i, msg := range testMessages {
		err := producer.Publish(ctx, []byte(msg))
		require.NoError(t, err)
		t.Logf("发送消息 %d: %s", i+1, msg)
	}

	// 等待消息被消费
	time.Sleep(2 * time.Second)
	t.Log("RabbitMQ 普通消息测试完成")
}

// TestRabbitMQDelayPublish 测试RabbitMQ延迟消息
func TestRabbitMQDelayPublish(t *testing.T) {
	cfg := facade.Config.Get()
	if !cfg.Rabbitmq.Enabled {
		t.Skip("RabbitMQ not enabled, skipping test")
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxkey.TraceIdKey, "test-delay-trace-id")

	// 等待消费者启动
	time.Sleep(500 * time.Millisecond)

	// 获取延迟生产者
	producer := facade.Queue.Producer("rabbitmq_delay_demo")
	if producer == nil {
		t.Skip("RabbitMQ delay producer not registered")
	}
	defer func() {
		err := producer.Close()
		if err != nil {
			t.Logf("关闭生产者失败: %v", err)
		}
	}()

	// 发送延迟消息(延迟5秒)
	startTime := time.Now()
	t.Logf("开始发送延迟消息: %v", startTime)

	err := producer.Publish(ctx, []byte(`{"orderId":111, "delay":5000, "message":"delay message 1"}`))
	require.NoError(t, err)

	err = producer.Publish(ctx, []byte(`{"orderId":222, "delay":5000, "message":"delay message 2"}`))
	require.NoError(t, err)

	t.Log("延迟消息已发送,等待6秒后消费...")

	// 等待消息被消费
	time.Sleep(6 * time.Second)

	endTime := time.Now()
	t.Logf("测试完成，耗时: %v", endTime.Sub(startTime))
}

// TestRabbitMQMultipleMessages 测试批量消息
func TestRabbitMQMultipleMessages(t *testing.T) {
	cfg := facade.Config.Get()
	if !cfg.Rabbitmq.Enabled {
		t.Skip("RabbitMQ not enabled, skipping test")
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxkey.TraceIdKey, "test-batch")

	// 等待消费者启动
	time.Sleep(500 * time.Millisecond)

	producer := facade.Queue.Producer("rabbitmq_demo")
	if producer == nil {
		t.Skip("RabbitMQ producer not registered")
	}
	defer func() {
		err := producer.Close()
		if err != nil {
			t.Errorf("关闭生产者失败: %s", err.Error())
		}
	}()

	// 批量发送20条消息
	batchSize := 20
	startTime := time.Now()

	for i := 1; i <= batchSize; i++ {
		msg := []byte(`{"orderId":` + string(rune('0'+i%10)) + `, "seq":` + string(rune('0'+i%10)) + `, "message":"batch message ` + string(rune('0'+i%10)) + `"}`)
		err := producer.Publish(ctx, msg)
		require.NoError(t, err)
	}

	t.Logf("批量发送 %d 条消息完成，耗时: %v", batchSize, time.Since(startTime))

	// 等待消息被消费
	time.Sleep(3 * time.Second)
	t.Log("批量测试完成")
}

// TestRabbitMQConsumerStatus 测试消费者状态
func TestRabbitMQConsumerStatus(t *testing.T) {
	cfg := facade.Config.Get()
	if !cfg.Rabbitmq.Enabled {
		t.Skip("RabbitMQ not enabled, skipping test")
	}

	// 获取所有消费者
	consumers := facade.Queue.GetAllConsumers()
	if len(consumers) == 0 {
		t.Skip("No consumers registered")
	}

	for _, c := range consumers {
		t.Logf("消费者: %s, 状态: %s, 启用: %v", c.Name(), c.Status(), c.Enabled(cfg))
	}

	// 获取消费者统计
	stats := facade.Queue.GetAllConsumerStats()
	for _, s := range stats {
		t.Logf("统计 - 消费者: %s, 状态: %s, 启用: %v", s.Name, s.Status, s.Enabled)
	}
}
