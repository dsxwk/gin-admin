package provider

import (
	"gin/app/queue"
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/logger"
	sq "gin/pkg/serviceprovider/queue"
	"fmt"
)

// QueueProvider 队列服务提供者
type QueueProvider struct {
	manager *sq.Manager
	log     *logger.Logger
}

// Name 服务提供者名称
func (p *QueueProvider) Name() string {
	return serviceprovider.ServiceQueue
}

// Register 注册队列服务和队列定义
func (p *QueueProvider) Register(app *container.Container) {
	manager, err := sq.NewManager(app.Config, queue.Consumers(), queue.Producers())
	if err != nil {
		flag.Errorf("队列服务注册失败: %v", err)
		return
	}

	app.SetQueue(manager)
	p.manager = manager
	flag.Infof(fmt.Sprintf("已注册 %d 个消费者, %d 个生产者", len(manager.Consumers()), len(manager.Producers())))
}

// Boot 启动服务
func (p *QueueProvider) Boot(app *container.Container) {
	cfg := app.Config()
	log := app.Log()
	p.log = log
	if cfg == nil {
		return
	}

	if p.manager == nil {
		return
	}
	if err := p.manager.Start(); err != nil {
		log.Error(fmt.Sprintf("队列服务启动失败: %v", err))
	}
}

// Runners 后台运行任务
func (p *QueueProvider) Runners() []serviceprovider.Runner {
	manager := p.manager
	log := p.log
	return []serviceprovider.Runner{
		serviceprovider.NewShutdownRunner("queue_shutdown", func() error {
			if manager == nil {
				return nil
			}
			if err := manager.Stop(); err != nil {
				log.Error(fmt.Sprintf("队列服务关闭失败: %v", err))
			}
			return nil
		}),
	}
}

// Dependencies 依赖服务
func (p *QueueProvider) Dependencies() []string {
	return []string{
		serviceprovider.ServiceConfig,
		serviceprovider.ServiceLog,
		serviceprovider.ServiceCache,
		serviceprovider.ServiceEvent,
	}
}
