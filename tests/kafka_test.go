package tests

import (
	"context"
	"gin/app/facade"
	"gin/app/queue/consumer"
	_ "gin/app/queue/producer"
	"gin/common/ctxkey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TestKafkaPublish Kafka普通消息发布和消费
func TestKafkaPublish(t *testing.T) {
	cfg := facade.Config()
	if !cfg.Queue.Kafka.Enabled {
		t.Skip("Kafka未启用, 跳过测试")
	}

	ctx := context.WithValue(t.Context(), ctxkey.TraceIdKey, "test-kafka-publish")

	time.Sleep(500 * time.Millisecond)

	producer := facade.Queue().Producer("kafka_demo")
	if producer == nil {
		t.Skip("Kafka生产者未注册")
	}

	testCases := []struct {
		name    string
		payload consumer.KafkaDemoPayload
	}{
		{"order_1", consumer.KafkaDemoPayload{Name: "kafka_test_1"}},
		{"order_2", consumer.KafkaDemoPayload{Name: "kafka_test_2"}},
		{"order_3", consumer.KafkaDemoPayload{Name: "kafka_test_3"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := producer.Publish(ctx, tc.payload)
			require.NoError(t, err, "发送消息失败: %s", tc.name)
			t.Logf("发送成功: %+v", tc.payload)
		})
	}

	time.Sleep(2 * time.Second)
	t.Log("Kafka普通消息测试完成")
}

// TestKafkaDelayPublish Kafka延迟消息发布和消费
func TestKafkaDelayPublish(t *testing.T) {
	cfg := facade.Config()
	if !cfg.Queue.Kafka.Enabled {
		t.Skip("Kafka未启用, 跳过测试")
	}

	ctx := context.WithValue(t.Context(), ctxkey.TraceIdKey, "test-kafka-delay")

	time.Sleep(500 * time.Millisecond)

	producer := facade.Queue().Producer("kafka_delay_demo")
	if producer == nil {
		t.Skip("Kafka延迟生产者未注册")
	}

	startTime := time.Now()
	t.Logf("开始发送延迟消息: %v", startTime)

	err := producer.Publish(ctx, consumer.KafkaDelayDemoPayload{Name: "kafka_delay_1"})
	require.NoError(t, err, "发送延迟消息1失败")

	err = producer.Publish(ctx, consumer.KafkaDelayDemoPayload{Name: "kafka_delay_2"})
	require.NoError(t, err, "发送延迟消息2失败")

	t.Log("延迟消息已发送, 等待6秒后消费...")

	time.Sleep(6 * time.Second)
	t.Logf("测试完成, 耗时: %v", time.Since(startTime))
}

// TestKafkaConsumerStatus 测试Kafka消费者状态查询
func TestKafkaConsumerStatus(t *testing.T) {
	cfg := facade.Config()
	if !cfg.Queue.Kafka.Enabled {
		t.Skip("Kafka未启用, 跳过测试")
	}

	consumers := facade.Queue().GetAllConsumers()
	if len(consumers) == 0 {
		t.Skip("未注册消费者")
	}

	for _, c := range consumers {
		name := c.Name()
		status := c.Status()
		enabled := c.Enabled(cfg)

		t.Run(name, func(t *testing.T) {
			t.Logf("消费者: %s, 状态: %s, 启用: %v", name, status, enabled)
			assert.NotEmpty(t, name, "消费者名称不能为空")

		})
	}
}
