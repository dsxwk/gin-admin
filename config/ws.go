package config

import "time"

// Ws WebSocket服务配置
type Ws struct {
	Enabled         bool          `mapstructure:"enabled" yaml:"enabled"`                     // 是否启用WebSocket服务
	Path            string        `mapstructure:"path" yaml:"path"`                           // WebSocket升级地址
	Shards          int           `mapstructure:"shards" yaml:"shards"`                       // 连接分片数量
	MaxConnections  int64         `mapstructure:"max-connections" yaml:"max-connections"`     // 最大连接数
	WriteQueueSize  int           `mapstructure:"write-queue-size" yaml:"write-queue-size"`   // 单连接发送队列长度
	ReadBufferSize  int           `mapstructure:"read-buffer-size" yaml:"read-buffer-size"`   // 读缓冲区大小
	WriteBufferSize int           `mapstructure:"write-buffer-size" yaml:"write-buffer-size"` // 写缓冲区大小
	ReadLimit       int64         `mapstructure:"read-limit" yaml:"read-limit"`               // 单条消息最大字节数
	WriteWait       time.Duration `mapstructure:"write-wait" yaml:"write-wait"`               // 写入超时
	PongWait        time.Duration `mapstructure:"pong-wait" yaml:"pong-wait"`                 // Pong等待时间
	PingInterval    time.Duration `mapstructure:"ping-interval" yaml:"ping-interval"`         // Ping发送间隔
	CloseWait       time.Duration `mapstructure:"close-wait" yaml:"close-wait"`               // 关闭等待时间
	AllowedOrigins  []string      `mapstructure:"allowed-origins" yaml:"allowed-origins"`     // 允许的Origin,空为同源
	Compression     bool          `mapstructure:"compression" yaml:"compression"`             // 是否启用压缩
}
