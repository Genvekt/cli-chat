package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/IBM/sarama"

	"github.com/Genvekt/cli-chat/services/auth_producer/internal/config"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
)

var _ config.KafkaProducerConfig = (*kafkaProducerConfig)(nil)

type kafkaProducerConfig struct {
	brokers []string
}

// NewKafkaProducerConfig reads kafka producer config from env
func NewKafkaProducerConfig() (*kafkaProducerConfig, error) {
	brokersStr := os.Getenv(brokersEnvName)
	if len(brokersStr) == 0 {
		return nil, fmt.Errorf("environment variable %s has not been set", brokersEnvName)
	}

	brokers := strings.Split(brokersStr, ",")

	return &kafkaProducerConfig{
		brokers: brokers,
	}, nil
}

// Brokers returns list of urls of kafka brokers
func (cfg *kafkaProducerConfig) Brokers() []string {
	return cfg.brokers
}

// Config возвращает конфигурацию для sarama producer
func (cfg *kafkaProducerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	return config
}
