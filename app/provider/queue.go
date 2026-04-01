package provider

import (
	"context"
	"gin/app/facade"
	_ "gin/app/queue/kafka/consumer"
	_ "gin/app/queue/kafka/producer"
	_ "gin/app/queue/rabbitmq/consumer"
	_ "gin/app/queue/rabbitmq/producer"
	"gin/pkg"
	"gin/pkg/foundation"
	"gin/pkg/queue"
)

func init() {
	foundation.Register(&QueueProvider{})
}

// QueueProvider 队列服务提供者
// 管理所有消费者的生命周期
type QueueProvider struct {
	consumers []queue.Consumer
	producers []queue.Producer
}

// Name 服务提供者名称
func (p *QueueProvider) Name() string {
	return "queue"
}

// Register 注册服务到门面
func (p *QueueProvider) Register(app foundation.App) {
	// 注册队列门面
	facade.Register("queue", facade.Queue)
	p.consumers = queue.GetConsumerRegistry().GetAll()
	p.producers = queue.GetProducerRegistry().GetAll()
	facade.Log.Info(pkg.Sprintf("已注册 %d 个消费者, %d 个生产者", len(p.consumers), len(p.producers)))
}

// Boot 启动服务(只启动消费者,生产者按需使用)
func (p *QueueProvider) Boot(app foundation.App) {
	cfg := facade.Config.Get()
	log := facade.Log.Logger()

	if cfg == nil {
		return
	}

	// 启动所有启用的消费者
	for _, consumer := range p.consumers {
		if consumer.Enabled(cfg) {
			facade.Log.Info(pkg.Sprintf("启动消费者: %s", consumer.Name()))
			if err := consumer.Start(cfg, log); err != nil {
				facade.Log.Error(pkg.Sprintf("启动消费者 %s 失败: %v", consumer.Name(), err))
			}
		}
	}

	// 获取所有生产者引用
	p.producers = facade.Queue.GetAllProducers()
}

// Runners 后台运行任务
func (p *QueueProvider) Runners() []foundation.Runner {
	return []foundation.Runner{
		&queueShutdownRunner{
			consumers: p.consumers,
			producers: p.producers,
		},
	}
}

// Dependencies 依赖服务
func (p *QueueProvider) Dependencies() []string {
	return []string{"config", "log"}
}

// queueShutdownRunner 队列关闭任务
type queueShutdownRunner struct {
	consumers []queue.Consumer
	producers []queue.Producer
}

// Run 运行等待任务
func (r *queueShutdownRunner) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// Stop 停止时关闭所有消费者和生产者
func (r *queueShutdownRunner) Stop() error {
	// 停止消费者
	for _, consumer := range r.consumers {
		if err := consumer.Stop(); err != nil {
			facade.Log.Error(pkg.Sprintf("停止消费者 %s 失败: %v", consumer.Name(), err))
		}
	}

	// 关闭生产者
	for _, producer := range r.producers {
		if err := producer.Close(); err != nil {
			facade.Log.Error(pkg.Sprintf("关闭生产者 %s 失败: %v", producer.Name(), err))
		}
	}

	facade.Log.Info("所有队列服务已关闭")
	return nil
}

// Name 任务名称
func (r *queueShutdownRunner) Name() string {
	return "queue_shutdown"
}
