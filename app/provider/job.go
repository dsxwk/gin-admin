package provider

import (
	"context"
	"gin/app/facade"
	_ "gin/app/job"
	"gin/common/flag"
	"gin/config"
	"gin/pkg"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/job"
	"gin/pkg/serviceprovider/logger"
	"gin/pkg/serviceprovider/queue"
	"time"

	"github.com/segmentio/kafka-go"
)

func init() {
	serviceprovider.Register(&JobProvider{})
}

// JobProvider 任务服务提供者
type JobProvider struct {
	consumers []queue.Consumer
	log       *logger.Logger
}

// Name 服务提供者名称
func (p *JobProvider) Name() string {
	return serviceprovider.ServiceJob
}

// Register 注册服务到容器
func (p *JobProvider) Register(app *container.Container) {
	app.Set(serviceprovider.ServiceJob, facade.NewJobFacade())
	flag.Infof(pkg.Sprintf("已注册 %d 个Job", job.Count()))
}

// Boot 启动服务
func (p *JobProvider) Boot(app *container.Container) {
	cfg := app.Get[*config.Config](serviceprovider.ServiceConfig)
	p.log = app.Get[*logger.Logger](serviceprovider.ServiceLog)
	if cfg == nil {
		return
	}

	// 收集所有注册job使用的connection
	connSet := make(map[string]bool)
	for _, jb := range job.GetAll() {
		c := jb.Connection()
		if c == "" {
			c = "redis"
		}
		connSet[c] = true
	}

	// 按connection创建任务消费者
	for conn := range connSet {
		switch conn {
		case "redis":
			p.consumers = append(p.consumers, job.NewRedisConsumer(&queue.RedisConsumer{
				Queue:        "job:queue",
				DelayedQueue: "job:queue:delayed",
				GetClient:    facade.RedisClient,
				Log:          p.log,
				Dual:         true,
			}))
		case "kafka":
			if cfg.Queue.Kafka.Enabled {
				driver := queue.NewKafka(cfg, p.log, facade.Event().Bus())
				driver.Reader = kafka.NewReader(kafka.ReaderConfig{
					Brokers:        cfg.Queue.Kafka.Brokers,
					Topic:          "job",
					GroupID:        "job_group",
					MinBytes:       1,
					MaxBytes:       10e6,
					StartOffset:    kafka.LastOffset,
					CommitInterval: 0,
					MaxWait:        5 * time.Second,
				})
				p.consumers = append(p.consumers, job.NewKafkaConsumer(&queue.KafkaConsumer{
					Kafka: driver,
					Topic: "job",
					Group: "job_group",
				}))
			} else {
				flag.Warningf("存在使用 kafka 连接的Job, 但 queue.kafka.enabled 为 false")
			}
		case "rabbitmq":
			if cfg.Queue.Rabbitmq.Enabled {
				driver, err := queue.NewRabbitMQ(cfg, p.log, facade.Event().Bus())
				if err != nil {
					p.log.Error(pkg.Sprintf("Job RabbitMQ 连接失败: %v", err))
					continue
				}
				p.consumers = append(p.consumers, job.NewRabbitmqConsumer(&queue.RabbitmqConsumer{
					Mq:       driver,
					Queue:    "job",
					Exchange: "job_exchange",
					Routing:  "job",
				}))
			} else {
				flag.Warningf("存在使用 rabbitmq 连接的Job, 但 queue.rabbitmq.enabled 为 false")
			}
		}
	}

	for _, consumer := range p.consumers {
		if err := consumer.Start(); err != nil {
			p.log.Error(pkg.Sprintf("Job消费者 %s 启动失败: %v", consumer.Name(), err))
			continue
		}
		flag.Infof("Job消费者已启动: %s", consumer.Name())
	}
}

// Runners 后台运行任务
func (p *JobProvider) Runners() []serviceprovider.Runner {
	return []serviceprovider.Runner{
		&jobShutdownRunner{
			consumers: p.consumers,
			log:       p.log,
		},
	}
}

// Dependencies 依赖服务
func (p *JobProvider) Dependencies() []string {
	return []string{
		serviceprovider.ServiceConfig,
		serviceprovider.ServiceLog,
		serviceprovider.ServiceEvent,
		serviceprovider.ServiceCache,
	}
}

// jobShutdownRunner 任务关闭器
type jobShutdownRunner struct {
	consumers []queue.Consumer
	log       *logger.Logger
}

// Run 运行等待
func (r *jobShutdownRunner) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// Stop 停止所有Worker
func (r *jobShutdownRunner) Stop() error {
	for _, consumer := range r.consumers {
		if err := consumer.Stop(); err != nil {
			r.log.Error(pkg.Sprintf("Job消费者 %s 关闭失败: %v", consumer.Name(), err))
		}
	}
	flag.Infof("所有Job消费者已关闭")
	return nil
}

// Name 任务名称
func (r *jobShutdownRunner) Name() string {
	return "job_shutdown"
}
