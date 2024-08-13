package producer

import (
	"log"

	"github.com/IBM/sarama"

	"github.com/Genvekt/cli-chat/libraries/kafka/pkg/kafka"
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

	log.Printf("message sent to partition %d with offset %d\n", partition, offset)

	return nil
}

func (p *kafkaProducer) Close() error {
	return p.producer.Close()
}
