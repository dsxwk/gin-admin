package ws

import "time"

// Options WebSocket运行配置
type Options struct {
	Shards          int           // 连接分片数量
	MaxConnections  int64         // 最大连接数
	WriteQueueSize  int           // 单连接发送队列长度
	ReadBufferSize  int           // 读缓冲区大小
	WriteBufferSize int           // 写缓冲区大小
	ReadLimit       int64         // 单条消息最大字节数
	WriteWait       time.Duration // 写入超时
	PongWait        time.Duration // Pong等待时间
	PingInterval    time.Duration // Ping发送间隔
	CloseWait       time.Duration // 关闭等待时间
	AllowedOrigins  []string      // 允许的Origin,空为同源
	Compression     bool          // 是否启用压缩
}

// DefaultOptions 默认WebSocket配置
func DefaultOptions() Options {
	return Options{
		Shards:          64,
		MaxConnections:  100000,
		WriteQueueSize:  32,
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		ReadLimit:       1 << 20,
		WriteWait:       10 * time.Second,
		PongWait:        60 * time.Second,
		PingInterval:    25 * time.Second,
		CloseWait:       5 * time.Second,
		Compression:     false,
	}
}

// normalize 标准化配置
func (o Options) normalize() Options {
	defaults := DefaultOptions()
	if o.Shards <= 0 {
		o.Shards = defaults.Shards
	}
	if o.MaxConnections <= 0 {
		o.MaxConnections = defaults.MaxConnections
	}
	if int64(o.Shards) > o.MaxConnections {
		o.Shards = int(o.MaxConnections)
	}
	if o.WriteQueueSize <= 0 {
		o.WriteQueueSize = defaults.WriteQueueSize
	}
	if o.ReadBufferSize <= 0 {
		o.ReadBufferSize = defaults.ReadBufferSize
	}
	if o.WriteBufferSize <= 0 {
		o.WriteBufferSize = defaults.WriteBufferSize
	}
	if o.ReadLimit <= 0 {
		o.ReadLimit = defaults.ReadLimit
	}
	if o.WriteWait <= 0 {
		o.WriteWait = defaults.WriteWait
	}
	if o.PongWait <= 0 {
		o.PongWait = defaults.PongWait
	}
	if o.PingInterval <= 0 || o.PingInterval >= o.PongWait {
		o.PingInterval = defaults.PingInterval
		if o.PingInterval >= o.PongWait {
			o.PingInterval = o.PongWait / 2
		}
	}
	if o.CloseWait <= 0 {
		o.CloseWait = defaults.CloseWait
	}

	return o
}
