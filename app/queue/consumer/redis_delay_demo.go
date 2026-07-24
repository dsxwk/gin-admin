package consumer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/config"
	"gin/pkg"
	"gin/pkg/serviceprovider/queue"
)

// RedisDelayDemoConsumer Redis延迟消费者
type RedisDelayDemoConsumer struct {
	*base.RedisConsumer
}

// RedisDelayDemoPayload 延迟消息体
type RedisDelayDemoPayload struct {
	Name string `json:"name"`
}

func (c *RedisDelayDemoConsumer) NewPayload() any {
	return &RedisDelayDemoPayload{}
}

func (c *RedisDelayDemoConsumer) Handle(payload any) error {
	data := payload.(*RedisDelayDemoPayload)
	facade.Log().Info(pkg.Sprintf("Redis Delay Received Msg: name=%s", data.Name))
	return nil
}

func NewRedisDelayDemoConsumer() *RedisDelayDemoConsumer {
	return &RedisDelayDemoConsumer{
		RedisConsumer: &base.RedisConsumer{
			Queue: "redis_delay_demo",
		},
	}
}

func (c *RedisDelayDemoConsumer) Name() string {
	return "redis_delay_demo"
}

func (c *RedisDelayDemoConsumer) Connection() string { return "redis" }

func (c *RedisDelayDemoConsumer) Retry() int { return 3 }

func (c *RedisDelayDemoConsumer) IsDelay() bool { return true }

func (c *RedisDelayDemoConsumer) Description() string {
	return "redis延迟队列消费者"
}

func (c *RedisDelayDemoConsumer) Start() error {
	c.RedisConsumer.Start(c)
	flag.Infof("Redis延迟消费者启动成功: %s", c.Name())
	return nil
}

func (c *RedisDelayDemoConsumer) Stop() error {
	return c.RedisConsumer.Stop()
}

func (c *RedisDelayDemoConsumer) Enabled(cfg *config.Config) bool {
	return true
}

func (c *RedisDelayDemoConsumer) Status() queue.ConsumerStatus {
	return c.RedisConsumer.Status()
}

func init() {
	queue.GetConsumerRegistry().Register(NewRedisDelayDemoConsumer())
}
