package consumer

import (
	"context"
	"errors"
	"strings"

	"github.com/IBM/sarama"

	"github.com/Genvekt/cli-chat/libraries/kafka/pkg/kafka"
	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
)

var _ kafka.Consumer[sarama.ConsumerMessage] = (*kafkaConsumer)(nil)

type Handler func(ctx context.Context, msg *sarama.ConsumerMessage) error

type kafkaConsumer struct {
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *GroupHandler
}

// NewConsumer creates consumer for kafka
func NewConsumer(
	consumerGroup sarama.ConsumerGroup,
	consumerGroupHandler *GroupHandler,
) *kafkaConsumer {
	return &kafkaConsumer{
		consumerGroup:        consumerGroup,
		consumerGroupHandler: consumerGroupHandler,
	}
}

func (c *kafkaConsumer) Consume(ctx context.Context, topic string, handler kafka.Handler[sarama.ConsumerMessage]) error {
	c.consumerGroupHandler.msgHandler = handler

	return c.consume(ctx, topic)
}

func (c *kafkaConsumer) Close() error {
	return c.consumerGroup.Close()
}

func (c *kafkaConsumer) consume(ctx context.Context, topicName string) error {
	for {
		err := c.consumerGroup.Consume(ctx, strings.Split(topicName, ","), c.consumerGroupHandler)
		if err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}

			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		logger.Info("rebalancing kafka consumer..")
	}
}
