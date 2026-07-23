package work

import (
	"context"
	"errors"
	"gin/app/facade"
	"gin/common/flag"
	"gin/pkg"
	"gin/pkg/serviceprovider/job"
	"github.com/goccy/go-json"
	"github.com/segmentio/kafka-go"
	"time"
)

// KafkaWorker Kafka任务消费者
type KafkaWorker struct {
	reader *kafka.Reader
	ctx    context.Context
	cancel context.CancelFunc
}

const JobKafkaTopic = "job"
const JobKafkaGroup = "job_group"

// NewKafkaWorker 创建KafkaWorker
func NewKafkaWorker() *KafkaWorker {
	cfg := facade.Config()
	if cfg == nil {
		return nil
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Queue.Kafka.Brokers,
		Topic:          JobKafkaTopic,
		GroupID:        JobKafkaGroup,
		MinBytes:       1,
		MaxBytes:       10e6,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 0,
		MaxWait:        5 * time.Second,
	})

	return &KafkaWorker{reader: reader}
}

func (w *KafkaWorker) Start() error {
	w.ctx, w.cancel = context.WithCancel(context.Background())
	go func() {
		flag.Infof("Job Kafka Worker 已启动")
		defer flag.Infof("Job Kafka Worker 已停止")
		for {
			select {
			case <-w.ctx.Done():
				return
			default:
				w.consume()
			}
		}
	}()
	return nil
}

func (w *KafkaWorker) Stop() error {
	if w.cancel != nil {
		w.cancel()
	}
	if w.reader != nil {
		return w.reader.Close()
	}
	return nil
}

func (w *KafkaWorker) consume() {
	msg, err := w.reader.FetchMessage(w.ctx)
	if err != nil {
		if !errors.Is(err, context.Canceled) {
			time.Sleep(time.Second)
		}
		return
	}
	var jm JobMessage
	if err = json.Unmarshal(msg.Value, &jm); err != nil {
		facade.Log().Error(pkg.Sprintf("Kafka Job 消息解析失败: %v", err))
		return
	}
	w.handleMessage(jm)
	_ = w.reader.CommitMessages(w.ctx, msg)
}

func (w *KafkaWorker) handleMessage(msg JobMessage) {
	j := job.Get(msg.JobName)
	if j == nil {
		facade.Log().Error(pkg.Sprintf("Job [%s] 未注册", msg.JobName))
		return
	}

	// 延迟执行
	if msg.RunAt > 0 {
		now := time.Now().UnixMilli()
		if msg.RunAt > now {
			time.Sleep(time.Millisecond * time.Duration(msg.RunAt-now))
		}
	}

	payload := j.NewPayload()
	if err := json.Unmarshal(msg.Payload, payload); err != nil {
		facade.Log().Error(pkg.Sprintf("Job [%s] payload解析失败: %v", msg.JobName, err))
		return
	}

	retry := j.Retry()
	delay := j.Delay()
	var err error
	for attempt := 0; attempt < retry || attempt == 0; attempt++ {
		err = j.Handle(payload)
		if err == nil {
			return
		}
		facade.Log().Error(pkg.Sprintf("Job [%s] 处理失败(attempt %d/%d): %v", msg.JobName, attempt+1, retry, err))
		if delay > 0 {
			time.Sleep(time.Millisecond * time.Duration(delay))
		}
	}
}

// JobMessage 消息格式
type JobMessage struct {
	JobName string          `json:"jobName"`
	Payload json.RawMessage `json:"payload"`
	RunAt   int64           `json:"runAt,omitempty"`
}
