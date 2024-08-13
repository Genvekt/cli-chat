package app

import (
	"log"

	"github.com/IBM/sarama"

	"github.com/Genvekt/cli-chat/libraries/closer/pkg/closer"
	"github.com/Genvekt/cli-chat/libraries/kafka/pkg/kafka"
	"github.com/Genvekt/cli-chat/libraries/kafka/pkg/kafka/producer"
	"github.com/Genvekt/cli-chat/services/auth_producer/internal/config/env"

	userImpl "github.com/Genvekt/cli-chat/services/auth_producer/internal/api/user"
	"github.com/Genvekt/cli-chat/services/auth_producer/internal/config"
	"github.com/Genvekt/cli-chat/services/auth_producer/internal/service"
	producerService "github.com/Genvekt/cli-chat/services/auth_producer/internal/service/producer"
)

// ServiceProvider initialises and stores various dependencies as singletons
type ServiceProvider struct {
	httpConfig config.HTTPConfig

	kafkaProducerConfig config.KafkaProducerConfig
	syncProducer        sarama.SyncProducer
	kafkaProducer       kafka.Producer

	userCreatorConfig  config.UserCreatorConfig
	userCreatorService service.UserCreatorService

	userAPI *userImpl.Service
}

func newServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

// HTTPConfig loads config related to http service of this application
func (s *ServiceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cong, err := env.NewHTTPConfigEnv()
		if err != nil {
			log.Fatalf("error loading http config: %v", err)
		}

		s.httpConfig = cong
	}

	return s.httpConfig
}

// KafkaProducerConfig loads config related to kafka producer
func (s *ServiceProvider) KafkaProducerConfig() config.KafkaProducerConfig {
	if s.kafkaProducerConfig == nil {
		conf, err := env.NewKafkaProducerConfig()
		if err != nil {
			log.Fatalf("Error loading kafka producer config: %v", err)
		}

		s.kafkaProducerConfig = conf
	}

	return s.kafkaProducerConfig
}

// SyncProducer initialises sarama sync producer
func (s *ServiceProvider) SyncProducer() sarama.SyncProducer {
	if s.syncProducer == nil {
		syncProducer, err := sarama.NewSyncProducer(
			s.KafkaProducerConfig().Brokers(),
			s.KafkaProducerConfig().Config(),
		)
		if err != nil {
			log.Fatalf("Error creating kafka sync producer: %v", err)
		}

		s.syncProducer = syncProducer
	}

	return s.syncProducer
}

// KafkaProducer initialises kafka producer
func (s *ServiceProvider) KafkaProducer() kafka.Producer {
	if s.kafkaProducer == nil {
		s.kafkaProducer = producer.NewProducer(s.SyncProducer())

		closer.Add(s.kafkaProducer.Close)
	}

	return s.kafkaProducer
}

// UserCreatorConfig loads config related to user creator
func (s *ServiceProvider) UserCreatorConfig() config.UserCreatorConfig {
	if s.userCreatorConfig == nil {
		conf, err := env.NewUserCreatorConfigEnv()
		if err != nil {
			log.Fatalf("Error loading user creator config: %v", err)
		}

		s.userCreatorConfig = conf
	}

	return s.userCreatorConfig
}

// UserCreatorService provides user creator
func (s *ServiceProvider) UserCreatorService() service.UserCreatorService {
	if s.userCreatorService == nil {
		s.userCreatorService = producerService.NewUserCreatorService(
			s.UserCreatorConfig(),
			s.KafkaProducer(),
		)
	}

	return s.userCreatorService
}

// UserAPI inits user api
func (s *ServiceProvider) UserAPI() *userImpl.Service {
	if s.userAPI == nil {
		s.userAPI = userImpl.NewService(s.UserCreatorService())
	}

	return s.userAPI
}
