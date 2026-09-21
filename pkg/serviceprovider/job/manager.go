package job

import (
	"context"
	"encoding/json"
	"fmt"
	"gin/common/ctxkey"
	"gin/common/flag"
	"gin/config"
	"gin/pkg/serviceprovider/debugger"
	"gin/pkg/serviceprovider/eventbus"
	"gin/pkg/serviceprovider/logger"
	"gin/pkg/serviceprovider/queue"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/kafka-go"
)

const (
	// JobKafkaTopic 任务Kafka主题
	JobKafkaTopic = "job"
	// JobRabbitmqQueue 任务RabbitMQ队列
	JobRabbitmqQueue = "job"
	// JobRabbitmqExchange 任务RabbitMQ交换机
	JobRabbitmqExchange = "job_exchange"
	// JobRabbitmqRouting 任务RabbitMQ路由
	JobRabbitmqRouting = "job"
	// JobRabbitmqDelayExchange 任务延迟交换机
	JobRabbitmqDelayExchange = "job_delay_exchange"
	// JobRabbitmqDelayQueue 任务延迟队列
	JobRabbitmqDelayQueue = "job_delay_queue"
)

// Manager 任务管理器
type Manager struct {
	config      func() *config.Config
	log         *logger.Logger
	redisClient func() *redis.Client
	bus         *eventbus.Bus
	jobs        map[string]Job

	mu           sync.Mutex
	rabbitmqConn *amqp091.Connection
	rabbitmqCh   *amqp091.Channel
	kafkaWriter  *kafka.Writer
	connInit     sync.Once
}

// NewManager 创建任务管理器
func NewManager(
	config func() *config.Config,
	log *logger.Logger,
	redisClient func() *redis.Client,
	bus *eventbus.Bus,
	jobs []Job,
) (*Manager, error) {
	manager := &Manager{
		config:      config,
		log:         log,
		redisClient: redisClient,
		bus:         bus,
		jobs:        make(map[string]Job, len(jobs)),
	}

	for _, item := range jobs {
		if item == nil {
			continue
		}

		name := item.Name()
		if name == "" {
			return nil, fmt.Errorf("job名称不能为空")
		}
		if _, exists := manager.jobs[name]; exists {
			return nil, fmt.Errorf("job [%s] 重复注册", name)
		}

		manager.jobs[name] = item
	}

	return manager, nil
}

// Config 获取配置
func (m *Manager) Config() *config.Config {
	if m == nil || m.config == nil {
		return nil
	}
	return m.config()
}

// Dispatch 投递任务
func (m *Manager) Dispatch(ctx context.Context, jobName string, payload any) error {
	if m == nil {
		return fmt.Errorf("job manager 未初始化")
	}

	item := m.jobs[jobName]
	if item == nil {
		return fmt.Errorf("job [%s] 未注册", jobName)
	}

	start := time.Now()
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("job [%s] payload序列化失败: %w", jobName, err)
	}

	connection := item.Connection()
	if connection == "" {
		connection = "redis"
	}
	delayMs := item.Delay()

	var dispatchErr error
	switch connection {
	case "sync":
		dispatchErr = m.dispatchSync(ctx, jobName, payloadBytes)
	case "redis":
		dispatchErr = m.dispatchRedis(ctx, jobName, payloadBytes, delayMs)
	case "kafka":
		dispatchErr = m.dispatchKafka(ctx, jobName, payloadBytes, delayMs)
	case "rabbitmq":
		dispatchErr = m.dispatchRabbitmq(ctx, jobName, payloadBytes, delayMs)
	default:
		dispatchErr = fmt.Errorf("job [%s] 不支持的连接: %s", jobName, connection)
	}

	if m != nil && m.bus != nil {
		m.bus.Publish(ctx, debugger.TopicJob, debugger.JobEvent{
			TraceID:    ctxkey.TraceID(ctx),
			Name:       jobName,
			Connection: connection,
			Payload:    string(payloadBytes),
			Ms:         float64(time.Since(start).Nanoseconds()) / 1e6,
		})
	}

	if dispatchErr != nil {
		flag.Errorf(dispatchErr.Error())
	}

	return dispatchErr
}

// dispatchSync 同步投递
func (m *Manager) dispatchSync(ctx context.Context, jobName string, payloadBytes []byte) error {
	return Execute(ctx, m.Job(jobName), NewMessage(jobName, payloadBytes, 0))
}

// dispatchRedis Redis投递
func (m *Manager) dispatchRedis(ctx context.Context, jobName string, payloadBytes []byte, delayMs int64) error {
	client := m.redis()
	if client == nil {
		return fmt.Errorf("redis 客户端未初始化")
	}

	delay := time.Duration(delayMs) * time.Millisecond
	message := NewMessage(jobName, payloadBytes, delay)
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return queue.PublishRedisMessage(
		ctx,
		client,
		"job:queue",
		"job:queue:delayed",
		messageBytes,
		delay,
	)
}

// dispatchKafka Kafka投递
func (m *Manager) dispatchKafka(ctx context.Context, jobName string, payloadBytes []byte, delayMs int64) error {
	m.initConnections()

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.kafkaWriter == nil {
		return fmt.Errorf("kafka writer 未初始化")
	}

	message := NewMessage(jobName, payloadBytes, time.Duration(delayMs)*time.Millisecond)
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return m.kafkaWriter.WriteMessages(ctx, kafka.Message{Value: messageBytes})
}

// dispatchRabbitmq RabbitMQ投递
func (m *Manager) dispatchRabbitmq(ctx context.Context, jobName string, payloadBytes []byte, delayMs int64) error {
	m.initConnections()

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.rabbitmqCh == nil {
		return fmt.Errorf("rabbitmq channel 未初始化")
	}

	message := NewMessage(jobName, payloadBytes, 0)
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	if delayMs > 0 {
		return m.rabbitmqCh.PublishWithContext(
			ctx,
			JobRabbitmqDelayExchange,
			JobRabbitmqRouting,
			false,
			false,
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        messageBytes,
				Expiration:  fmt.Sprintf("%d", delayMs),
			},
		)
	}

	return m.rabbitmqCh.PublishWithContext(
		ctx,
		JobRabbitmqExchange,
		JobRabbitmqRouting,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        messageBytes,
		},
	)
}

// initConnections 初始化持久连接
func (m *Manager) initConnections() {
	m.connInit.Do(func() {
		cfg := m.Config()
		if cfg == nil {
			return
		}

		if cfg.Queue.Rabbitmq.Enabled && cfg.Queue.Rabbitmq.Url != "" {
			conn, err := amqp091.Dial(cfg.Queue.Rabbitmq.Url)
			if err != nil {
				m.logError("Job RabbitMQ 连接失败: " + err.Error())
			} else {
				ch, err := conn.Channel()
				if err != nil {
					m.logError("Job RabbitMQ Channel 创建失败: " + err.Error())
					_ = conn.Close()
				} else {
					_ = ch.ExchangeDeclare(JobRabbitmqExchange, "direct", true, false, false, false, nil)
					_, _ = ch.QueueDeclare(JobRabbitmqQueue, true, false, false, false, nil)
					_ = ch.QueueBind(JobRabbitmqQueue, JobRabbitmqRouting, JobRabbitmqExchange, false, nil)
					_ = ch.ExchangeDeclare(JobRabbitmqDelayExchange, "direct", true, false, false, false, nil)
					_, _ = ch.QueueDeclare(JobRabbitmqDelayQueue, true, false, false, false, amqp091.Table{
						"x-dead-letter-exchange":    JobRabbitmqExchange,
						"x-dead-letter-routing-key": JobRabbitmqRouting,
					})
					_ = ch.QueueBind(JobRabbitmqDelayQueue, JobRabbitmqRouting, JobRabbitmqDelayExchange, false, nil)
					m.rabbitmqConn = conn
					m.rabbitmqCh = ch
				}
			}
		}

		if cfg.Queue.Kafka.Enabled && len(cfg.Queue.Kafka.Brokers) > 0 {
			m.kafkaWriter = &kafka.Writer{
				Addr:     kafka.TCP(cfg.Queue.Kafka.Brokers...),
				Topic:    JobKafkaTopic,
				Balancer: &kafka.LeastBytes{},
			}
		}
	})
}

// redis 获取Redis客户端
func (m *Manager) redis() *redis.Client {
	if m == nil || m.redisClient == nil {
		return nil
	}
	return m.redisClient()
}

// logError 记录错误日志
func (m *Manager) logError(message string) {
	if m != nil && m.log != nil {
		m.log.Error(message)
	}
}

// JobStats 任务统计
type JobStats struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Connection  string `json:"connection"`
}

// GetAllJobs 获取所有任务
func (m *Manager) GetAllJobs() []JobStats {
	items := m.Jobs()
	stats := make([]JobStats, 0, len(items))
	for _, item := range items {
		connection := item.Connection()
		if connection == "" {
			connection = "redis"
		}
		stats = append(stats, JobStats{
			Name:        item.Name(),
			Description: item.Description(),
			Connection:  connection,
		})
	}
	return stats
}

// Jobs 获取所有任务
func (m *Manager) Jobs() []Job {
	if m == nil {
		return nil
	}

	items := make([]Job, 0, len(m.jobs))
	for _, item := range m.jobs {
		items = append(items, item)
	}

	return items
}

// Job 获取指定任务
func (m *Manager) Job(name string) Job {
	if m == nil {
		return nil
	}
	return m.jobs[name]
}

// Count 统计待处理任务
func (m *Manager) Count(ctx context.Context) (int64, error) {
	client := m.redis()
	if client == nil {
		return 0, fmt.Errorf("redis 客户端未初始化")
	}
	listCount, _ := client.LLen(ctx, "job:queue").Result()
	zCount, _ := client.ZCard(ctx, "job:queue:delayed").Result()
	return listCount + zCount, nil
}

// Clear 清除待处理任务
func (m *Manager) Clear(ctx context.Context) error {
	client := m.redis()
	if client == nil {
		return fmt.Errorf("redis 客户端未初始化")
	}
	client.Del(ctx, "job:queue")
	client.Del(ctx, "job:queue:delayed")
	return nil
}
