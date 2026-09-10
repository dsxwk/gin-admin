package debugger

const (
	TopicSql      = "debug:sql"
	TopicCache    = "debug:cache"
	TopicHttp     = "debug:http"
	TopicMq       = "debug:mq"
	TopicGrpc     = "debug:grpc"
	TopicListener = "debug:listener"
	TopicJob      = "debug:job"
	TopicEs       = "debug:es"
)

// SqlEvent Sql事件
type SqlEvent struct {
	TraceId string // tranceId
	Sql     string
	Rows    int64
	Ms      float64
}

// CacheEvent 缓存事件
type CacheEvent struct {
	TraceId string // tranceId
	Driver  string
	Name    string
	Cmd     string
	Args    any
	Ms      float64
}

// HttpEvent Http事件
type HttpEvent struct {
	TraceId  string // tranceId
	Url      string
	Method   string
	Header   map[string]string
	Body     any
	Status   int
	Response any
	Ms       float64
}

// MqEvent 消息队列事件
type MqEvent struct {
	TraceId string         // tranceId
	Driver  string         // kafka rabbitmq redis-stream
	Topic   string         // topic queue stream
	Message string         // 消息内容
	Key     string         // 用于Kafka
	Group   string         // 消费组
	Ms      float64        // 耗时ms
	Extra   map[string]any // 扩展信息
}

// GrpcEvent grpc事件
type GrpcEvent struct {
	TraceId  string
	Method   string
	Request  any
	Response any
	Code     string
	Ms       float64
}

// EsEvent ES事件
type EsEvent struct {
	TraceId  string  // traceId
	Action   string  // 操作
	Index    string  // 索引
	Request  any     // 请求
	Response any     // 响应
	Code     string  // 状态
	Ms       float64 // 耗时
}

// ListenerEvent 监听事件
type ListenerEvent struct {
	TraceId     string // tranceId
	Name        string // 监听名称
	Description string // 监听描述
	Data        any    // 监听数据
}

// JobEvent Job事件
type JobEvent struct {
	TraceId    string  `json:"traceId"`    // traceId
	Name       string  `json:"name"`       // Job名称
	Connection string  `json:"connection"` // 连接类型
	Payload    string  `json:"payload"`    // 消息内容
	Ms         float64 `json:"ms"`         // 耗时ms
}
