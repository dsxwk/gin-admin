package job

import (
	"encoding/json"
	"time"
)

// Message 任务消息
type Message struct {
	JobName string          `json:"jobName"`
	Payload json.RawMessage `json:"payload"`
	RunAt   int64           `json:"runAt,omitempty"`
}

// NewMessage 创建任务消息
func NewMessage(jobName string, payload []byte, delay time.Duration) Message {
	message := Message{
		JobName: jobName,
		Payload: payload,
	}
	if delay > 0 {
		message.RunAt = time.Now().Add(delay).UnixMilli()
	}
	return message
}
