package provider

import (
	"context"
	"gin/app/facade"
	_ "gin/app/job"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/job"
	"gin/pkg/serviceprovider/job/work"
)

func init() {
	serviceprovider.Register(&JobProvider{})
}

// JobProvider 任务服务提供者
type JobProvider struct {
	redisWorker    *work.RedisWorker
	kafkaWorker    *work.KafkaWorker
	rabbitmqWorker *work.RabbitmqWorker
}

// Name 服务提供者名称
func (p *JobProvider) Name() string {
	return "job"
}

// Register 注册服务到门面
func (p *JobProvider) Register(app serviceprovider.App) {
	facade.Register("job", facade.Job())
	flag.Infof(pkg.Sprintf("已注册 %d 个Job", job.Count()))
}

// Boot 启动服务
func (p *JobProvider) Boot(app serviceprovider.App) {
	cfg := facade.Config()
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

	// 按connection启动对应的Worker
	for conn := range connSet {
		switch conn {
		case "redis":
			p.redisWorker = work.NewRedisWorker()
			if p.redisWorker != nil {
				flag.Infof("启动 Job Redis Worker")
				if err := p.redisWorker.Start(); err != nil {
					facade.Log().Error(pkg.Sprintf("Job Redis Worker 启动失败: %v", err))
				}
			}
		case "kafka":
			if cfg.Queue.Kafka.Enabled {
				p.kafkaWorker = work.NewKafkaWorker()
				if p.kafkaWorker != nil {
					flag.Infof("启动 Job Kafka Worker")
					if err := p.kafkaWorker.Start(); err != nil {
						facade.Log().Error(pkg.Sprintf("Job Kafka Worker 启动失败: %v", err))
					}
				}
			} else {
				flag.Warningf("存在使用 kafka 连接的Job, 但 queue.kafka.enabled 为 false")
			}
		case "rabbitmq":
			if cfg.Queue.Rabbitmq.Enabled {
				p.rabbitmqWorker = work.NewRabbitmqWorker()
				if p.rabbitmqWorker != nil {
					flag.Infof("启动 Job RabbitMQ Worker")
					if err := p.rabbitmqWorker.Start(); err != nil {
						facade.Log().Error(pkg.Sprintf("Job RabbitMQ Worker 启动失败: %v", err))
					}
				}
			} else {
				flag.Warningf("存在使用 rabbitmq 连接的Job, 但 queue.rabbitmq.enabled 为 false")
			}
		}
	}
}

// Runners 后台运行任务
func (p *JobProvider) Runners() []serviceprovider.Runner {
	return []serviceprovider.Runner{
		&jobShutdownRunner{
			redisWorker:    p.redisWorker,
			kafkaWorker:    p.kafkaWorker,
			rabbitmqWorker: p.rabbitmqWorker,
		},
	}
}

// Dependencies 依赖服务
func (p *JobProvider) Dependencies() []string {
	return []string{"config", "log"}
}

// jobShutdownRunner 任务关闭器
type jobShutdownRunner struct {
	redisWorker    *work.RedisWorker
	kafkaWorker    *work.KafkaWorker
	rabbitmqWorker *work.RabbitmqWorker
}

// Run 运行等待
func (r *jobShutdownRunner) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// Stop 停止所有Worker
func (r *jobShutdownRunner) Stop() error {
	if r.redisWorker != nil {
		if err := r.redisWorker.Stop(); err != nil {
			facade.Log().Error(pkg.Sprintf("Job Redis Worker 关闭失败: %v", err))
		}
	}
	if r.kafkaWorker != nil {
		if err := r.kafkaWorker.Stop(); err != nil {
			facade.Log().Error(pkg.Sprintf("Job Kafka Worker 关闭失败: %v", err))
		}
	}
	if r.rabbitmqWorker != nil {
		if err := r.rabbitmqWorker.Stop(); err != nil {
			facade.Log().Error(pkg.Sprintf("Job RabbitMQ Worker 关闭失败: %v", err))
		}
	}
	flag.Infof("所有Job Worker已关闭")
	return nil
}

// Name 任务名称
func (r *jobShutdownRunner) Name() string {
	return "job_shutdown"
}
