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
	"sync"
	"time"
)

const JobKafkaTopic = "job"
const JobKafkaGroup = "job_group"
const kafkaWorkerPoolSize = 3

// KafkaWorker Kafka任务消费者(worker pool模式,消费和延迟处理分离)
type KafkaWorker struct {
	reader *kafka.Reader
	ctx    context.CancelFunc
	wg     sync.WaitGroup
}

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

// Start 启动消费者(消费goroutine + worker pool)
func (w *KafkaWorker) Start() error {
	baseCtx, cancel := context.WithCancel(context.Background())
	w.ctx = cancel

	msgCh := make(chan kafka.Message, 10)

	// 启动worker pool, 处理业务(允许time.Sleep延迟不阻塞消费)
	for i := 0; i < kafkaWorkerPoolSize; i++ {
		w.wg.Add(1)
		go w.processWorker(baseCtx, msgCh)
	}

	// 启动消费goroutine, 只负责拉取消息
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		defer close(msgCh)
		flag.Infof("Job Kafka Worker 已启动")
		defer flag.Infof("Job Kafka Worker 已停止")
		for {
			msg, err := w.reader.FetchMessage(baseCtx)
			if err != nil {
				if !errors.Is(err, context.Canceled) {
					time.Sleep(time.Second)
					continue
				}
				return
			}
			select {
			case <-baseCtx.Done():
				return
			case msgCh <- msg:
			}
		}
	}()

	return nil
}

// Stop 停止消费者
func (w *KafkaWorker) Stop() error {
	if w.ctx != nil {
		w.ctx()
	}
	w.wg.Wait()
	if w.reader != nil {
		return w.reader.Close()
	}
	return nil
}

// processWorker worker池, 处理消息(延迟通过time.Sleep实现)
func (w *KafkaWorker) processWorker(ctx context.Context, msgCh <-chan kafka.Message) {
	defer w.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-msgCh:
			if !ok {
				return
			}
			var jm JobMessage
			if err := json.Unmarshal(msg.Value, &jm); err != nil {
				facade.Log().Error(pkg.Sprintf("Kafka Job 消息解析失败: %v", err))
				w.commit(ctx, msg)
				continue
			}
			w.handleMessage(ctx, jm)
			w.commit(ctx, msg)
		}
	}
}

// commit 提交消息偏移
func (w *KafkaWorker) commit(ctx context.Context, msg kafka.Message) {
	if err := w.reader.CommitMessages(ctx, msg); err != nil {
		facade.Log().Error(pkg.Sprintf("Kafka Job Commit失败: %v", err))
	}
}

// handleMessage 处理消息(延迟在worker池内执行,不阻塞消费)
func (w *KafkaWorker) handleMessage(ctx context.Context, msg JobMessage) {
	j := job.Get(msg.JobName)
	if j == nil {
		facade.Log().Error(pkg.Sprintf("Job [%s] 未注册", msg.JobName))
		return
	}

	// 延迟执行(Kafka无原生延迟, worker池内time.Sleep不阻塞其他消息消费)
	if msg.RunAt > 0 {
		now := time.Now().UnixMilli()
		if msg.RunAt > now {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Millisecond * time.Duration(msg.RunAt-now)):
			}
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
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Millisecond * time.Duration(delay)):
			}
		}
	}
}

// JobMessage 消息格式
type JobMessage struct {
	JobName string          `json:"jobName"`
	Payload json.RawMessage `json:"payload"`
	RunAt   int64           `json:"runAt,omitempty"`
}
