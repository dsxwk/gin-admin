package facade

import (
	"context"
	"encoding/json"
	"fmt"
	"gin/common/ctxkey"
	"gin/common/flag"
	"gin/pkg/container"
	"gin/pkg/serviceprovider"
	"gin/pkg/serviceprovider/debugger"
	jsjob "gin/pkg/serviceprovider/job"
	"gin/pkg/serviceprovider/queue"
	"sync"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/kafka-go"
)

const JobKafkaTopic = "job"
const JobRabbitmqQueue = "job"
const JobRabbitmqExchange = "job_exchange"
const JobRabbitmqRouting = "job"
const JobRabbitmqDelayExchange = "job_delay_exchange"
const JobRabbitmqDelayQueue = "job_delay_queue"

// Job 任务门面实例
func Job() *JobFacade {
	return container.Default().Get[*JobFacade](serviceprovider.ServiceJob)
}

// NewJobFacade 创建任务门面
func NewJobFacade() *JobFacade {
	return &JobFacade{}
}

type JobFacade struct {
	mu           sync.Mutex
	rabbitmqConn *amqp091.Connection
	rabbitmqCh   *amqp091.Channel
	kafkaWriter  *kafka.Writer
	connInit     sync.Once
}

// initConnections 初始化持久连接(仅执行一次)
func (j *JobFacade) initConnections() {
	j.connInit.Do(func() {
		cfg := Config()
		if cfg == nil {
			return
		}
		if cfg.Queue.Rabbitmq.Enabled && cfg.Queue.Rabbitmq.Url != "" {
			conn, err := amqp091.Dial(cfg.Queue.Rabbitmq.Url)
			if err != nil {
				Log().Error("Job RabbitMQ 连接失败: " + err.Error())
			} else {
				ch, err := conn.Channel()
				if err != nil {
					Log().Error("Job RabbitMQ Channel 创建失败: " + err.Error())
					_ = conn.Close()
				} else {
					_ = ch.ExchangeDeclare(JobRabbitmqExchange, "direct", true, false, false, false, nil)
					_, _ = ch.QueueDeclare(JobRabbitmqQueue, true, false, false, false, nil)
					_ = ch.QueueBind(JobRabbitmqQueue, JobRabbitmqRouting, JobRabbitmqExchange, false, nil)
					// 延迟队列: 使用DLX模式, 消息到期后自动路由到普通队列
					_ = ch.ExchangeDeclare(JobRabbitmqDelayExchange, "direct", true, false, false, false, nil)
					_, _ = ch.QueueDeclare(JobRabbitmqDelayQueue, true, false, false, false, amqp091.Table{
						"x-dead-letter-exchange":    JobRabbitmqExchange,
						"x-dead-letter-routing-key": JobRabbitmqRouting,
					})
					_ = ch.QueueBind(JobRabbitmqDelayQueue, JobRabbitmqRouting, JobRabbitmqDelayExchange, false, nil)
					j.rabbitmqConn = conn
					j.rabbitmqCh = ch
				}
			}
		}
		if cfg.Queue.Kafka.Enabled && len(cfg.Queue.Kafka.Brokers) > 0 {
			j.kafkaWriter = &kafka.Writer{
				Addr:     kafka.TCP(cfg.Queue.Kafka.Brokers...),
				Topic:    JobKafkaTopic,
				Balancer: &kafka.LeastBytes{},
			}
		}
	})
}

func (j *JobFacade) Dispatch(ctx context.Context, jobName string, payload any) error {
	jb := jsjob.Get(jobName)
	if jb == nil {
		return fmt.Errorf("job [%s] 未注册", jobName)
	}

	start := time.Now()
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("job [%s] payload序列化失败: %w", jobName, err)
	}

	conn := jb.Connection()
	if conn == "" {
		conn = "redis"
	}
	delayMs := jb.Delay()

	var dispatchErr error
	switch conn {
	case "sync":
		dispatchErr = j.dispatchSync(ctx, jobName, payloadBytes)
	case "redis":
		dispatchErr = j.dispatchRedis(ctx, jobName, payloadBytes, delayMs)
	case "kafka":
		dispatchErr = j.dispatchKafka(ctx, jobName, payloadBytes, delayMs)
	case "rabbitmq":
		dispatchErr = j.dispatchRabbitmq(ctx, jobName, payloadBytes, delayMs)
	default:
		dispatchErr = fmt.Errorf("job [%s] 不支持的连接: %s", jobName, conn)
	}

	Event().Bus().Publish(debugger.TopicJob, debugger.JobEvent{
		TraceID:    ctxkey.GetTraceId(ctx),
		Name:       jobName,
		Connection: conn,
		Payload:    string(payloadBytes),
		Ms:         float64(time.Since(start).Nanoseconds()) / 1e6,
	})

	if dispatchErr != nil {
		flag.Errorf(dispatchErr.Error())
	}

	return dispatchErr
}

func (j *JobFacade) dispatchSync(ctx context.Context, jobName string, payloadBytes []byte) error {
	return jsjob.Execute(ctx, jsjob.NewMessage(jobName, payloadBytes, 0))
}

func (j *JobFacade) dispatchRedis(ctx context.Context, jobName string, payloadBytes []byte, delayMs int64) error {
	client := Redis().Client()
	if client == nil {
		return fmt.Errorf("redis 客户端未初始化")
	}

	if delayMs > 0 {
		msg := jsjob.NewMessage(jobName, payloadBytes, time.Duration(delayMs)*time.Millisecond)
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		return queue.PublishRedisMessage(
			ctx,
			client,
			"job:queue",
			"job:queue:delayed",
			msgBytes,
			time.Duration(delayMs)*time.Millisecond,
		)
	}

	msg := jsjob.NewMessage(jobName, payloadBytes, 0)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return queue.PublishRedisMessage(ctx, client, "job:queue", "job:queue:delayed", msgBytes, 0)
}

func (j *JobFacade) dispatchKafka(ctx context.Context, jobName string, payloadBytes []byte, delayMs int64) error {
	j.initConnections()

	j.mu.Lock()
	defer j.mu.Unlock()

	if j.kafkaWriter == nil {
		return fmt.Errorf("kafka writer 未初始化")
	}
	msg := jsjob.NewMessage(jobName, payloadBytes, time.Duration(delayMs)*time.Millisecond)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return j.kafkaWriter.WriteMessages(ctx, kafka.Message{Value: msgBytes})
}

func (j *JobFacade) dispatchRabbitmq(ctx context.Context, jobName string, payloadBytes []byte, delayMs int64) error {
	j.initConnections()

	j.mu.Lock()
	defer j.mu.Unlock()

	if j.rabbitmqCh == nil {
		return fmt.Errorf("rabbitmq channel 未初始化")
	}
	// 延迟消息: 发送到延迟exchange, 通过TTL+DLX自动到期后路由到普通队列
	if delayMs > 0 {
		msg := jsjob.NewMessage(jobName, payloadBytes, 0)
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		return j.rabbitmqCh.PublishWithContext(
			ctx,
			JobRabbitmqDelayExchange,
			JobRabbitmqRouting,
			false,
			false,
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        msgBytes,
				Expiration:  fmt.Sprintf("%d", delayMs),
			},
		)
	}
	// 普通消息: 直接发送到普通exchange
	msg := jsjob.NewMessage(jobName, payloadBytes, 0)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return j.rabbitmqCh.PublishWithContext(
		ctx,
		JobRabbitmqExchange,
		JobRabbitmqRouting,
		false,
		false,
		amqp091.Publishing{ContentType: "application/json", Body: msgBytes},
	)
}

type JobStats struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Connection  string `json:"connection"`
}

func (j *JobFacade) GetAllJobs() []JobStats {
	jobs := jsjob.GetAll()
	stats := make([]JobStats, 0, len(jobs))
	for _, jb := range jobs {
		conn := jb.Connection()
		if conn == "" {
			conn = "redis"
		}
		stats = append(stats, JobStats{
			Name:        jb.Name(),
			Description: jb.Description(),
			Connection:  conn,
		})
	}
	return stats
}

func (j *JobFacade) Count(ctx context.Context) (int64, error) {
	client := Redis().Client()
	if client == nil {
		return 0, fmt.Errorf("redis 客户端未初始化")
	}
	listCount, _ := client.LLen(ctx, "job:queue").Result()
	zCount, _ := client.ZCard(ctx, "job:queue:delayed").Result()
	return listCount + zCount, nil
}

func (j *JobFacade) Clear(ctx context.Context) error {
	client := Redis().Client()
	if client == nil {
		return fmt.Errorf("redis 客户端未初始化")
	}
	client.Del(ctx, "job:queue")
	client.Del(ctx, "job:queue:delayed")
	return nil
}
