package consumer

import (
	"gin/app/facade"
	"gin/common/base"
	"gin/common/flag"
	"gin/config"
	"gin/pkg"
	"gin/pkg/serviceprovider/queue"
)

// RedisDemoConsumer Redis普通消费者
type RedisDemoConsumer struct {
	*base.RedisConsumer
}

// RedisDemoPayload 消息体
type RedisDemoPayload struct {
	Name string `json:"name"`
}

func (c *RedisDemoConsumer) NewPayload() any {
	return &RedisDemoPayload{}
}

func (c *RedisDemoConsumer) Handle(payload any) error {
	data := payload.(*RedisDemoPayload)
	facade.Log().Info(pkg.Sprintf("Redis Received Msg: name=%s", data.Name))
	return nil
}

func NewRedisDemoConsumer() *RedisDemoConsumer {
	return &RedisDemoConsumer{
		RedisConsumer: &base.RedisConsumer{
			Queue: "redis_demo",
		},
	}
}

func (c *RedisDemoConsumer) Name() string {
	return "redis_demo"
}

func (c *RedisDemoConsumer) Connection() string { return "redis" }

func (c *RedisDemoConsumer) Retry() int { return 3 }

func (c *RedisDemoConsumer) IsDelay() bool { return false }

func (c *RedisDemoConsumer) Description() string {
	return "redis普通队列消费者"
}

func (c *RedisDemoConsumer) Start() error {
	c.RedisConsumer.Start(c)
	flag.Infof("Redis消费者启动成功: %s", c.Name())
	return nil
}

func (c *RedisDemoConsumer) Stop() error {
	return c.RedisConsumer.Stop()
}

func (c *RedisDemoConsumer) Enabled(cfg *config.Config) bool {
	return true
}

func (c *RedisDemoConsumer) Status() queue.ConsumerStatus {
	return c.RedisConsumer.Status()
}

func init() {
	queue.GetConsumerRegistry().Register(NewRedisDemoConsumer())
}
