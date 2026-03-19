package base

import (
	"gin/pkg/logger"
	"gin/pkg/queue"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

type RabbitmqConsumer struct {
	Mq           *RabbitMq
	Queue        string
	Exchange     string
	Routing      string
	IsDelayQueue bool
	Retry        int
}

func (c *RabbitmqConsumer) Start(h queue.Handler) {
	go func() {
		for {
			if c.Mq == nil || c.Mq.Channel == nil {
				time.Sleep(time.Second)
				continue
			}

			args := amqp091.Table{}
			exchangeType := "direct"
			if c.IsDelayQueue {
				exchangeType = "x-delayed-message"
				args["x-delayed-type"] = "direct"
			}

			if err := c.Mq.Channel.ExchangeDeclare(
				c.Exchange,
				exchangeType,
				true,
				false,
				false,
				false,
				args,
			); err != nil {
				logger.NewLogger().Error("[RabbitMq] ExchangeDeclare error: " + err.Error())
				time.Sleep(time.Second)
				continue
			}

			if _, err := c.Mq.Channel.QueueDeclare(
				c.Queue,
				true,
				false,
				false,
				false,
				nil,
			); err != nil {
				logger.NewLogger().Error("[RabbitMq] QueueDeclare error: " + err.Error())
				time.Sleep(time.Second)
				continue
			}

			if err := c.Mq.Channel.QueueBind(
				c.Queue,
				c.Routing,
				c.Exchange,
				false,
				nil,
			); err != nil {
				logger.NewLogger().Error("[RabbitMq] QueueBind error: " + err.Error())
				time.Sleep(time.Second)
				continue
			}

			msgs, err := c.Mq.Channel.Consume(
				c.Queue,
				"",
				false,
				false,
				false,
				false,
				nil,
			)
			if err != nil {
				logger.NewLogger().Error("[RabbitMq] Consume error: " + err.Error())
				time.Sleep(time.Second)
				continue
			}

			for msg := range msgs {
				go func(msg amqp091.Delivery) {
					retry := 0
					for {
						err = h.Handle(string(msg.Body))
						if err == nil {
							if ackErr := msg.Ack(false); ackErr != nil {
								logger.NewLogger().Error("[RabbitMq] Ack error: " + ackErr.Error())
							}
							break
						}

						retry++
						if retry >= c.Retry {
							logger.NewLogger().Error("[RabbitMq] Retry failed: " + string(msg.Body))
							if ackErr := msg.Ack(false); ackErr != nil {
								logger.NewLogger().Error("[RabbitMq] Ack error: " + ackErr.Error())
							}
							break
						}
						time.Sleep(time.Second)
					}
				}(msg)
			}

			time.Sleep(time.Second)
		}
	}()
}
