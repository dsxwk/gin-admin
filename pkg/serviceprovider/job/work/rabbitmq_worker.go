package work

import (
	"encoding/json"
	"gin/app/facade"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/serviceprovider/job"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

// RabbitmqWorker RabbitMQ任务消费者
type RabbitmqWorker struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	stopCh  chan struct{}
}

const (
	JobRabbitmqQueue    = "job"
	JobRabbitmqExchange = "job_exchange"
	JobRabbitmqRouting  = "job"
)

func NewRabbitmqWorker() *RabbitmqWorker {
	cfg := facade.Config()
	if cfg == nil {
		return nil
	}
	conn, err := amqp091.Dial(cfg.Queue.Rabbitmq.Url)
	if err != nil {
		facade.Log().Error(pkg.Sprintf("RabbitMQ Job 连接失败: %v", err))
		return nil
	}
	ch, err := conn.Channel()
	if err != nil {
		facade.Log().Error(pkg.Sprintf("RabbitMQ Job Channel失败: %v", err))
		_ = conn.Close()
		return nil
	}
	return &RabbitmqWorker{conn: conn, channel: ch, stopCh: make(chan struct{})}
}

func (w *RabbitmqWorker) Start() error {
	if w.channel == nil {
		return nil
	}
	_ = w.channel.ExchangeDeclare(JobRabbitmqExchange, "direct", true, false, false, false, nil)
	_, _ = w.channel.QueueDeclare(JobRabbitmqQueue, true, false, false, false, nil)
	_ = w.channel.QueueBind(JobRabbitmqQueue, JobRabbitmqRouting, JobRabbitmqExchange, false, nil)
	msgs, err := w.channel.Consume(JobRabbitmqQueue, "", false, false, false, false, nil)
	if err != nil {
		facade.Log().Error(pkg.Sprintf("RabbitMQ Job Consume error: %v", err))
		return err
	}
	go func() {
		flag.Infof("Job RabbitMQ Worker 已启动")
		defer flag.Infof("Job RabbitMQ Worker 已停止")
		for {
			select {
			case <-w.stopCh:
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}
				w.handleMessage(msg)
			}
		}
	}()
	return nil
}

func (w *RabbitmqWorker) Stop() error {
	close(w.stopCh)
	if w.channel != nil {
		_ = w.channel.Close()
	}
	if w.conn != nil {
		return w.conn.Close()
	}
	return nil
}

func (w *RabbitmqWorker) handleMessage(msg amqp091.Delivery) {
	var jm JobMessage
	if err := json.Unmarshal(msg.Body, &jm); err != nil {
		facade.Log().Error(pkg.Sprintf("RabbitMQ Job 消息解析失败: %v", err))
		_ = msg.Ack(false)
		return
	}
	j := job.Get(jm.JobName)
	if j == nil {
		facade.Log().Error(pkg.Sprintf("Job [%s] 未注册", jm.JobName))
		_ = msg.Ack(false)
		return
	}

	if jm.RunAt > 0 {
		now := time.Now().UnixMilli()
		if jm.RunAt > now {
			time.Sleep(time.Millisecond * time.Duration(jm.RunAt-now))
		}
	}

	payload := j.NewPayload()
	if err := json.Unmarshal(jm.Payload, payload); err != nil {
		facade.Log().Error(pkg.Sprintf("Job [%s] payload解析失败: %v", jm.JobName, err))
		_ = msg.Ack(false)
		return
	}
	retry := j.Retry()
	delay := j.Delay()
	var err error
	for attempt := 0; attempt < retry || attempt == 0; attempt++ {
		err = j.Handle(payload)
		if err == nil {
			_ = msg.Ack(false)
			return
		}
		facade.Log().Error(pkg.Sprintf("Job [%s] 处理失败(attempt %d/%d): %v", jm.JobName, attempt+1, retry, err))
		if delay > 0 {
			time.Sleep(time.Millisecond * time.Duration(delay))
		}
	}
	_ = msg.Ack(false)
}
