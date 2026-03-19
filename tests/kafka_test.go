package tests

import (
	"context"
	"gin/app/queue/kafka/consumer"
	"gin/app/queue/kafka/producer"
	"gin/common/ctxkey"
	"gin/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKafkaPublish(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(context.Background(), ctxkey.TraceIdKey, "test-trace-id")

	if config.NewConfig().Kafka.Enabled {
		consumer.NewKafkaDemoConsumer()
		pdc := producer.NewKafkaDemoProducer()
		defer func() {
			require.NoError(t, pdc.Close())
		}()

		err := pdc.Publish(ctx, []byte(`{"orderId":111}`))
		require.NoError(t, err)
	}
}
