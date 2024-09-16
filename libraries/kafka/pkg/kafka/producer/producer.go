package producer

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/Genvekt/cli-chat/libraries/kafka/pkg/kafka"
	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
)

var _ kafka.Producer[sarama.ProducerMessage] = (*kafkaProducer)(nil)

type kafkaProducer struct {
	producer sarama.SyncProducer
}

// NewProducer creates producer for kafka
func NewProducer(producer sarama.SyncProducer) *kafkaProducer {
	return &kafkaProducer{
		producer: producer,
	}
}

func (p *kafkaProducer) SendMessage(msg *sarama.ProducerMessage) error {
	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	logger.Debug("kafka message sent",
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
		zap.String("topic", msg.Topic),
	)

	return nil
}

func (p *kafkaProducer) Close() error {
	return p.producer.Close()
}
