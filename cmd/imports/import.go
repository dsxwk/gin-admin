package imports

import (
	_ "gin/app/command"
	_ "gin/app/listener"
	_ "gin/app/queue/kafka/consumer"
	_ "gin/app/queue/kafka/producer"
	_ "gin/app/queue/rabbitmq/consumer"
	_ "gin/app/queue/rabbitmq/producer"
	_ "gin/pkg/cli/db"
	_ "gin/pkg/cli/event"
	_ "gin/pkg/cli/make"
	_ "gin/pkg/cli/queue"
	_ "gin/pkg/cli/route"
)
