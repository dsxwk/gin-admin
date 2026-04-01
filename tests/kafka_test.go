package tests

import (
	"context"
	"gin/app/facade"
	"gin/common/ctxkey"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestKafkaPublish 测试Kafka普通消息发布和消费
func TestKafkaPublish(t *testing.T) {
	cfg := facade.Config.Get()
	if !cfg.Kafka.Enabled {
		t.Skip("Kafka not enabled, skipping test")
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxkey.TraceIdKey, "test-trace-id")

	// 等待消费者启动
	time.Sleep(500 * time.Millisecond)

	// 获取生产者
	producer := facade.Queue.Producer("kafka_demo")
	if producer == nil {
		t.Skip("Kafka producer not registered")
	}
	defer func() {
		err := producer.Close()
		if err != nil {
			t.Errorf("Failed to close Kafka producer: %v", err)
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
	t.Log("Kafka 普通消息测试完成")
}

// TestKafkaDelayPublish 测试Kafka延迟消息
func TestKafkaDelayPublish(t *testing.T) {
	cfg := facade.Config.Get()
	if !cfg.Kafka.Enabled {
		t.Skip("Kafka not enabled, skipping test")
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxkey.TraceIdKey, "test-delay-trace-id")

	// 等待消费者启动
	time.Sleep(500 * time.Millisecond)

	// 获取延迟生产者
	producer := facade.Queue.Producer("kafka_delay_demo")
	if producer == nil {
		t.Skip("Kafka delay producer not registered")
	}
	defer func() {
		err := producer.Close()
		if err != nil {
			t.Errorf("Kafka delay producer close error: %s", err.Error())
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
	t.Logf("测试完成,耗时: %v", endTime.Sub(startTime))
}

// TestKafkaConsumerStatus 测试消费者状态
func TestKafkaConsumerStatus(t *testing.T) {
	cfg := facade.Config.Get()
	if !cfg.Kafka.Enabled {
		t.Skip("Kafka not enabled, skipping test")
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
