package facade

import (
	"context"
	"encoding/json"
	"fmt"
	"gin/common/ctxkey"
	"gin/common/flag"
	"gin/pkg/serviceprovider/debugger"
	jsjob "gin/pkg/serviceprovider/job"
	"github.com/go-redis/redis/v8"
	"github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

const JobKafkaTopic = "job"
const JobRabbitmqQueue = "job"
const JobRabbitmqExchange = "job_exchange"
const JobRabbitmqRouting = "job"

var (
	jobOnce   sync.Once
	jobFacade *JobFacade
)

// Job 任务门面实例(单例)
func Job() *JobFacade {
	jobOnce.Do(func() {
		jobFacade = &JobFacade{}
	})
	return jobFacade
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
					conn.Close()
				} else {
					_ = ch.ExchangeDeclare(JobRabbitmqExchange, "direct", true, false, false, false, nil)
					_, _ = ch.QueueDeclare(JobRabbitmqQueue, true, false, false, false, nil)
					_ = ch.QueueBind(JobRabbitmqQueue, JobRabbitmqRouting, JobRabbitmqExchange, false, nil)
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

func getTraceId(ctx context.Context) string {
	if ctx == nil {
		return "unknown"
	}
	if id := ctx.Value(ctxkey.TraceIdKey); id != nil {
		if s, ok := id.(string); ok && s != "" {
			return s
		}
	}
	return "unknown"
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
		dispatchErr = j.dispatchSync(jb, payloadBytes)
	case "redis":
		dispatchErr = j.dispatchRedis(ctx, jobName, payloadBytes, delayMs)
	case "kafka":
		dispatchErr = j.dispatchKafka(ctx, jobName, payloadBytes, delayMs)
	case "rabbitmq":
		dispatchErr = j.dispatchRabbitmq(ctx, jobName, payloadBytes, delayMs)
	default:
		dispatchErr = fmt.Errorf("job [%s] 不支持的连接: %s", jobName, conn)
	}

	Message().Publish(debugger.TopicJob, debugger.JobEvent{
		TraceId:    getTraceId(ctx),
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

func (j *JobFacade) dispatchSync(jb jsjob.Job, payloadBytes []byte) error {
	p := jb.NewPayload()
	if err := json.Unmarshal(payloadBytes, p); err != nil {
		return err
	}
	return jb.Handle(p)
}

func (j *JobFacade) dispatchRedis(ctx context.Context, jobName string, payloadBytes []byte, delayMs int64) error {
	client := Cache("redis").Redis().Client()
	if client == nil {
		return fmt.Errorf("redis 客户端未初始化")
	}

	if delayMs > 0 {
		msg := DelayedMessage{
			JobName: jobName,
			Payload: payloadBytes,
		}
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		score := float64(time.Now().UnixMilli() + delayMs)
		return client.ZAdd(ctx, "job:queue:delayed", &redis.Z{Score: score, Member: string(msgBytes)}).Err()
	}

	msg := buildJobMessage(jobName, payloadBytes, 0)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return client.LPush(ctx, "job:queue", string(msgBytes)).Err()
}

func (j *JobFacade) dispatchKafka(ctx context.Context, jobName string, payloadBytes []byte, delayMs int64) error {
	j.initConnections()

	j.mu.Lock()
	defer j.mu.Unlock()

	if j.kafkaWriter == nil {
		return fmt.Errorf("kafka writer 未初始化")
	}
	msg := buildJobMessage(jobName, payloadBytes, delayMs)
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
	msg := buildJobMessage(jobName, payloadBytes, delayMs)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return j.rabbitmqCh.Publish(
		JobRabbitmqExchange,
		JobRabbitmqRouting,
		false,
		false,
		amqp091.Publishing{ContentType: "application/json", Body: msgBytes},
	)
}

func buildJobMessage(jobName string, payloadBytes []byte, delayMs int64) JobMessage {
	msg := JobMessage{
		JobName: jobName,
		Payload: payloadBytes,
	}
	if delayMs > 0 {
		msg.RunAt = time.Now().UnixMilli() + delayMs
	}
	return msg
}

type JobMessage struct {
	JobName string          `json:"jobName"`
	Payload json.RawMessage `json:"payload"`
	RunAt   int64           `json:"runAt,omitempty"`
}

type DelayedMessage struct {
	JobName string          `json:"jobName"`
	Payload json.RawMessage `json:"payload"`
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
	client := Cache("redis").Redis().Client()
	if client == nil {
		return 0, fmt.Errorf("redis 客户端未初始化")
	}
	listCount, _ := client.LLen(ctx, "job:queue").Result()
	zCount, _ := client.ZCard(ctx, "job:queue:delayed").Result()
	return listCount + zCount, nil
}

func (j *JobFacade) Clear(ctx context.Context) error {
	client := Cache("redis").Redis().Client()
	if client == nil {
		return fmt.Errorf("redis 客户端未初始化")
	}
	client.Del(ctx, "job:queue")
	client.Del(ctx, "job:queue:delayed")
	return nil
}
