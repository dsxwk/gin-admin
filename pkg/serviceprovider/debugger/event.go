package debugger

// SQLEvent SQL调试事件
type SQLEvent struct {
	TraceID string  `json:"traceId"`
	SQL     string  `json:"sql"`
	Rows    int64   `json:"rows"`
	Ms      float64 `json:"ms"`
}

// CacheEvent 缓存调试事件
type CacheEvent struct {
	TraceID string  `json:"traceId"`
	Driver  string  `json:"driver"`
	Name    string  `json:"name"`
	Cmd     string  `json:"cmd"`
	Args    any     `json:"args"`
	Ms      float64 `json:"ms"`
}

// HTTPEvent HTTP调试事件
type HTTPEvent struct {
	TraceID  string            `json:"traceId"`
	URL      string            `json:"url"`
	Method   string            `json:"method"`
	Header   map[string]string `json:"header"`
	Body     any               `json:"body"`
	Status   int               `json:"status"`
	Response any               `json:"response"`
	Ms       float64           `json:"ms"`
}

// MQEvent 消息队列调试事件
type MQEvent struct {
	TraceID string         `json:"traceId"`
	Driver  string         `json:"driver"`
	Topic   string         `json:"topic"`
	Message string         `json:"message"`
	Key     string         `json:"key"`
	Group   string         `json:"group"`
	Ms      float64        `json:"ms"`
	Extra   map[string]any `json:"extra"`
}

// GRPCEvent gRPC调试事件
type GRPCEvent struct {
	TraceID  string  `json:"traceId"`
	Method   string  `json:"method"`
	Request  any     `json:"request"`
	Response any     `json:"response"`
	Code     string  `json:"code"`
	Ms       float64 `json:"ms"`
}

// ListenerEvent 业务监听调试事件
type ListenerEvent struct {
	TraceID     string `json:"traceId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Data        any    `json:"data"`
}

// JobEvent Job调试事件
type JobEvent struct {
	TraceID    string  `json:"traceId"`
	Name       string  `json:"name"`
	Connection string  `json:"connection"`
	Payload    string  `json:"payload"`
	Ms         float64 `json:"ms"`
}

// ESEvent ES调试事件
type ESEvent struct {
	TraceID  string  `json:"traceId"`
	Action   string  `json:"action"`
	Index    string  `json:"index"`
	Request  any     `json:"request"`
	Response any     `json:"response"`
	Code     string  `json:"code"`
	Ms       float64 `json:"ms"`
}
