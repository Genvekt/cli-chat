package producer

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"

	"github.com/Genvekt/cli-chat/libraries/kafka/pkg/kafka"
)

var _ kafka.Producer = (*kafkaProducer)(nil)

type kafkaProducer struct {
	producer sarama.SyncProducer
}

// NewProducer creates producer for kafka
func NewProducer(producer sarama.SyncProducer) *kafkaProducer {
	return &kafkaProducer{
		producer: producer,
	}
}

func (p *kafkaProducer) SendMessage(iMsg interface{}) error {
	msg, ok := iMsg.(*sarama.ProducerMessage)
	if !ok {
		return fmt.Errorf("unexpected message type: want %T, got %T", &sarama.ProducerMessage{}, iMsg)
	}

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
