package consumer

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/IBM/sarama"

	"github.com/Genvekt/cli-chat/libraries/kafka/pkg/kafka"
)

var _ kafka.Consumer = (*kafkaConsumer)(nil)

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

func (c *kafkaConsumer) Consume(ctx context.Context, topic string, handler kafka.Handler) error {
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

		log.Printf("rebalancing...\n")
	}
}
