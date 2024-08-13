package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/IBM/sarama"

	"github.com/Genvekt/cli-chat/services/auth/internal/config"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
	groupIDEnvName = "KAFKA_GROUP_ID"
)

var _ config.KafkaConsumerConfig = (*kafkaConsumerConfig)(nil)

type kafkaConsumerConfig struct {
	brokers []string
	groupID string
}

// NewKafkaConsumerConfig reads eonfig of kafka consumer from env
func NewKafkaConsumerConfig() (*kafkaConsumerConfig, error) {
	brokersStr := os.Getenv(brokersEnvName)
	if len(brokersStr) == 0 {
		return nil, fmt.Errorf("environment variable %s has not been set", brokersEnvName)
	}

	brokers := strings.Split(brokersStr, ",")

	groupID := os.Getenv(groupIDEnvName)
	if len(groupID) == 0 {
		return nil, fmt.Errorf("environment variable %s has not been set", groupIDEnvName)
	}

	return &kafkaConsumerConfig{
		brokers: brokers,
		groupID: groupID,
	}, nil
}

// Brokers returns the list of url of kafka brokers
func (cfg *kafkaConsumerConfig) Brokers() []string {
	return cfg.brokers
}

// GroupID returns the id of kafka consumers group
func (cfg *kafkaConsumerConfig) GroupID() string {
	return cfg.groupID
}

// Config возвращает конфигурацию для sarama consumer
func (cfg *kafkaConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
