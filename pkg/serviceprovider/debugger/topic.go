package debugger

import "gin/pkg/serviceprovider/eventbus"

const (
	// TopicSQL SQL调试
	TopicSQL = "debug:sql"
	// TopicCache 缓存调试
	TopicCache = "debug:cache"
	// TopicHTTP HTTP调试
	TopicHTTP = "debug:http"
	// TopicMQ 消息队列调试
	TopicMQ = "debug:mq"
	// TopicGRPC gRPC调试
	TopicGRPC = "debug:grpc"
	// TopicListener 业务事件调试
	TopicListener = eventbus.TopicEvent
	// TopicJob Job调试
	TopicJob = "debug:job"
	// TopicES ES调试
	TopicES = "debug:es"
)
