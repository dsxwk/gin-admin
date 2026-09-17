package provider

import (
	"context"
	"gin/app/facade"
	"gin/common/flag"
	"gin/config"
	"gin/pkg"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/logger"
	"gin/pkg/serviceprovider/queue"
)

func init() {
	serviceprovider.Register(&QueueProvider{})
}

// QueueProvider 队列服务提供者
type QueueProvider struct {
	consumers []queue.Consumer
	producers []queue.Producer
	log       *logger.Logger
}

// Name 服务提供者名称
func (p *QueueProvider) Name() string {
	return serviceprovider.ServiceQueue
}

// Register 注册队列服务和队列定义
func (p *QueueProvider) Register(app *container.Container) {
	app.Set(serviceprovider.ServiceQueue, facade.NewQueueFacade())
	consumers := queue.GetConsumerRegistry()
	producers := queue.GetProducerRegistry()

	if err := consumers.RegisterFactories(); err != nil {
		flag.Errorf("%v", err)
	}
	if err := producers.RegisterFactories(); err != nil {
		flag.Errorf("%v", err)
	}

	p.consumers = consumers.GetAll()
	p.producers = producers.GetAll()
	flag.Infof(pkg.Sprintf("已注册 %d 个消费者, %d 个生产者", len(p.consumers), len(p.producers)))
}

// Boot 启动服务
func (p *QueueProvider) Boot(app *container.Container) {
	cfg := app.Get[*config.Config](serviceprovider.ServiceConfig)
	log := app.Get[*logger.Logger](serviceprovider.ServiceLog)
	p.log = log
	if cfg == nil {
		return
	}

	for _, consumer := range p.consumers {
		if consumer.Enabled(cfg) {
			flag.Infof("启动消费者: %s", consumer.Name())
			if err := consumer.Start(); err != nil {
				log.Error(pkg.Sprintf("启动消费者 %s 失败: %v", consumer.Name(), err))
			}
		}
	}

	p.producers = queue.GetProducerRegistry().GetAll()
}

// Runners 后台运行任务
func (p *QueueProvider) Runners() []serviceprovider.Runner {
	return []serviceprovider.Runner{
		&queueShutdownRunner{
			consumers: p.consumers,
			producers: p.producers,
			log:       p.log,
		},
	}
}

// Dependencies 依赖服务
func (p *QueueProvider) Dependencies() []string {
	return []string{serviceprovider.ServiceConfig, serviceprovider.ServiceLog, serviceprovider.ServiceEvent}
}

// queueShutdownRunner 队列关闭任务
type queueShutdownRunner struct {
	consumers []queue.Consumer
	producers []queue.Producer
	log       *logger.Logger
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
			r.log.Error(pkg.Sprintf("停止消费者 %s 失败: %v", consumer.Name(), err))
		}
	}

	// 关闭生产者
	for _, producer := range r.producers {
		if err := producer.Close(); err != nil {
			r.log.Error(pkg.Sprintf("关闭生产者 %s 失败: %v", producer.Name(), err))
		}
	}

	flag.Infof("所有队列服务已关闭")
	return nil
}

// Name 任务名称
func (r *queueShutdownRunner) Name() string {
	return "queue_shutdown"
}
