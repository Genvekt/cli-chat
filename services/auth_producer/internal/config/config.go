package config

import (
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

// Load reads .env file into environment vars
func Load(filePath string) error {
	if filePath == "" {
		// nothing to load
		return nil
	}

	err := godotenv.Load(filePath)
	if err != nil {
		return err
	}

	return nil
}

// HTTPConfig provides parameters related to HTTP server
type HTTPConfig interface {
	Address() string
}

// KafkaProducerConfig provides parameters related to kafka producer
type KafkaProducerConfig interface {
	Brokers() []string
	Config() *sarama.Config
}

// UserKafkaClientConfig provides parameters related to user kafka client
type UserKafkaClientConfig interface {
	Topic() string
}
