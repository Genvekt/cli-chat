package config

import (
	"time"

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

// GRPCConfig provides parameters related to GRPC server
type GRPCConfig interface {
	Address() string
	TLSCertFile() string
	TLSKeyFile() string
	IsTLSEnabled() bool
}

// HTTPConfig provides parameters related to HTTP server
type HTTPConfig interface {
	Address() string
}

// PostgresConfig provides parameters related to Postgres database
type PostgresConfig interface {
	DSN() string
}

// UserServiceConfig provides parameters related to user service
type UserServiceConfig interface {
	CacheTTL() time.Duration
	NoCache() bool
	UseCache() bool
}

// KafkaConsumerConfig provides parameters related to kafka consumer
type KafkaConsumerConfig interface {
	Brokers() []string
	GroupID() string
	Config() *sarama.Config
}

// UserSaverConfig provides parameters related to user saver service
type UserSaverConfig interface {
	Topic() string
}

// TokenProviderConfig provides parameters related to token generation
type TokenProviderConfig interface {
	Secret() []byte
	TTL() time.Duration
}

// JaegerTracingConfig provides envs related to jaeger
type JaegerTracingConfig interface {
	ServiceName() string
	SamplerType() string
	SamplerParam() float64
	AgentAddress() string
}
