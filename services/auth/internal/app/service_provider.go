package app

import (
	"context"
	"log"

	"github.com/IBM/sarama"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/Genvekt/cli-chat/libraries/kafka/pkg/kafka"
	"github.com/Genvekt/cli-chat/libraries/kafka/pkg/kafka/consumer"
	"github.com/Genvekt/cli-chat/services/auth/internal/utils"
	"github.com/Genvekt/cli-chat/services/auth/internal/utils/hash"
	"github.com/Genvekt/cli-chat/services/auth/internal/utils/token"

	"github.com/Genvekt/cli-chat/libraries/closer/pkg/closer"

	"github.com/Genvekt/cli-chat/libraries/cache_client/pkg/cache"
	"github.com/Genvekt/cli-chat/libraries/cache_client/pkg/cache/redis"
	cacheConfig "github.com/Genvekt/cli-chat/libraries/cache_client/pkg/config"
	"github.com/Genvekt/cli-chat/libraries/db_client/pkg/db"
	"github.com/Genvekt/cli-chat/libraries/db_client/pkg/db/pg"
	"github.com/Genvekt/cli-chat/libraries/db_client/pkg/db/transaction"
	accessImpl "github.com/Genvekt/cli-chat/services/auth/internal/api/access"
	authImpl "github.com/Genvekt/cli-chat/services/auth/internal/api/auth"
	userImpl "github.com/Genvekt/cli-chat/services/auth/internal/api/user"
	"github.com/Genvekt/cli-chat/services/auth/internal/config"
	"github.com/Genvekt/cli-chat/services/auth/internal/config/env"
	"github.com/Genvekt/cli-chat/services/auth/internal/repository"
	accessRepository "github.com/Genvekt/cli-chat/services/auth/internal/repository/access/pg"
	userRepository "github.com/Genvekt/cli-chat/services/auth/internal/repository/user/pg"
	userCache "github.com/Genvekt/cli-chat/services/auth/internal/repository/user/redis"
	"github.com/Genvekt/cli-chat/services/auth/internal/service"
	accessService "github.com/Genvekt/cli-chat/services/auth/internal/service/access"
	authService "github.com/Genvekt/cli-chat/services/auth/internal/service/auth"
	consumerService "github.com/Genvekt/cli-chat/services/auth/internal/service/consumer"
	userService "github.com/Genvekt/cli-chat/services/auth/internal/service/user"
)

// ServiceProvider initialises and stores various dependencies as singletons
type ServiceProvider struct {
	gRPCConfig        config.GRPCConfig
	httpConfig        config.HTTPConfig
	swaggerConfig     config.HTTPConfig
	postgresConfig    config.PostgresConfig
	redisConfig       cacheConfig.RedisConfig
	userServiceConfig config.UserServiceConfig

	refreshTokenConfig config.TokenProviderConfig
	accessTokenConfig  config.TokenProviderConfig

	passwordHasher       utils.Hasher
	refreshTokenProvider utils.TokenProvider
	accessTokenProvider  utils.TokenProvider

	kafkaConsumerConfig  config.KafkaConsumerConfig
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *consumer.GroupHandler
	kafkaConsumer        kafka.Consumer[sarama.ConsumerMessage]

	dbClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	userRepo   repository.UserRepository
	accessRepo repository.AccessRepository

	userCache repository.UserCache

	userSaverConfig  config.UserSaverConfig
	userSaverService service.ConsumerService

	userService   service.UserService
	authService   service.AuthService
	accessService service.AccessService

	userImpl   *userImpl.Service
	authImpl   *authImpl.Service
	accessImpl *accessImpl.Service
}

func newServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

// GRPCConfig provides configuration of grpc server of this application
func (s *ServiceProvider) GRPCConfig() config.GRPCConfig {
	if s.gRPCConfig == nil {
		grpcConfig, err := env.NewGRPCConfigEnv()
		if err != nil {
			log.Fatalf("failed to load grpc config: %v", err)
		}

		s.gRPCConfig = grpcConfig
	}

	return s.gRPCConfig
}

// HTTPConfig provides configuration of http server of this application
func (s *ServiceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		httpConfig, err := env.NewHTTPConfigEnv()
		if err != nil {
			log.Fatalf("failed to load http config: %v", err)
		}

		s.httpConfig = httpConfig
	}

	return s.httpConfig
}

// SwaggerConfig provides configuration of swagger server of this application
func (s *ServiceProvider) SwaggerConfig() config.HTTPConfig {
	if s.swaggerConfig == nil {
		swaggerCongig, err := env.NewSwaggerConfigEnv()
		if err != nil {
			log.Fatalf("failed to load swagger config: %v", err)
		}

		s.swaggerConfig = swaggerCongig
	}

	return s.swaggerConfig
}

// PGConfig provides configuration parameters for postgres db
func (s *ServiceProvider) PGConfig() config.PostgresConfig {
	if s.postgresConfig == nil {
		postgresConfig, err := env.NewPostgresConfigEnv()
		if err != nil {
			log.Fatalf("failed to load postgres config: %v", err)
		}

		s.postgresConfig = postgresConfig
	}

	return s.postgresConfig
}

// RedisConfig provides configuration parameters for redis cache
func (s *ServiceProvider) RedisConfig() cacheConfig.RedisConfig {
	if s.redisConfig == nil {
		redisConfig, err := env.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to load redis config: %v", err)
		}

		s.redisConfig = redisConfig
	}

	return s.redisConfig
}

// KafkaConsumerConfig provides configuration parameters for kafka consumer
func (s *ServiceProvider) KafkaConsumerConfig() config.KafkaConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		kafkaConsumerConfig, err := env.NewKafkaConsumerConfig()
		if err != nil {
			log.Fatalf("failed to load kafka consumer config: %v", err)
		}

		s.kafkaConsumerConfig = kafkaConsumerConfig
	}

	return s.kafkaConsumerConfig
}

// RefreshTokenConfig provides configuration parameters for refresh token
func (s *ServiceProvider) RefreshTokenConfig() config.TokenProviderConfig {
	if s.refreshTokenConfig == nil {
		conf, err := env.NewRefreshTokenProviderConfig()
		if err != nil {
			log.Fatalf("failed to load refresh token config: %v", err)
		}

		s.refreshTokenConfig = conf
	}

	return s.refreshTokenConfig
}

// RefreshTokenProvider provides refresh token
func (s *ServiceProvider) RefreshTokenProvider() utils.TokenProvider {
	if s.refreshTokenProvider == nil {
		s.refreshTokenProvider = token.NewTokenProvider(s.RefreshTokenConfig())
	}

	return s.refreshTokenProvider
}

// AccessTokenConfig provides configuration parameters for access token
func (s *ServiceProvider) AccessTokenConfig() config.TokenProviderConfig {
	if s.accessTokenConfig == nil {
		conf, err := env.NewAccessTokenProviderConfig()
		if err != nil {
			log.Fatalf("failed to load access token config: %v", err)
		}

		s.accessTokenConfig = conf
	}

	return s.accessTokenConfig
}

// AccessTokenProvider provides access token
func (s *ServiceProvider) AccessTokenProvider() utils.TokenProvider {
	if s.accessTokenProvider == nil {
		s.accessTokenProvider = token.NewTokenProvider(s.AccessTokenConfig())
	}

	return s.accessTokenProvider
}

// PasswordHasher initialises password hasher
func (s *ServiceProvider) PasswordHasher() utils.Hasher {
	if s.passwordHasher == nil {
		s.passwordHasher = hash.NewHasher()
	}

	return s.passwordHasher
}

// DBClient provides DB client over postgres
func (s *ServiceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		pgClient, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to postgres: %v", err)
		}

		if err := pgClient.DB().Ping(ctx); err != nil {
			log.Fatalf("failed to ping postgres: %v", err)
		}

		closer.Add(func() error {
			return pgClient.Close()
		})

		s.dbClient = pgClient
	}

	return s.dbClient
}

// TxManager provides transaction manager
func (s *ServiceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// RedisPool creates connection to redis
func (s *ServiceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", s.RedisConfig().Address())
			},
		}

		closer.Add(s.redisPool.Close)
	}

	return s.redisPool
}

// RedisClient provides redis client dependency
func (s *ServiceProvider) RedisClient() cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	return s.redisClient
}

// UserCache provides user cache dependency
func (s *ServiceProvider) UserCache() repository.UserCache {
	if s.userCache == nil {
		s.userCache = userCache.NewUserCacheRedis(s.RedisClient())
	}

	return s.userCache
}

// UserRepo provides user repository dependency
func (s *ServiceProvider) UserRepo(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = userRepository.NewUserRepositoryPostgres(s.DBClient(ctx))
	}

	return s.userRepo
}

// AccessRepo provides access repository dependency
func (s *ServiceProvider) AccessRepo(ctx context.Context) repository.AccessRepository {
	if s.accessRepo == nil {
		s.accessRepo = accessRepository.NewAccessRepositoryPostgres(s.DBClient(ctx))
	}

	return s.accessRepo
}

// ConsumerGroup provides kafka consumer group
func (s *ServiceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			s.KafkaConsumerConfig().Brokers(),
			s.KafkaConsumerConfig().GroupID(),
			s.KafkaConsumerConfig().Config(),
		)
		if err != nil {
			log.Fatalf("failed to create consumer group: %v", err)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

// ConsumerGroupHandler provides handler for kafka consumer group
func (s *ServiceProvider) ConsumerGroupHandler() *consumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = consumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}

// KafkaConsumer initialises kafka consumer
func (s *ServiceProvider) KafkaConsumer() kafka.Consumer[sarama.ConsumerMessage] {
	if s.kafkaConsumer == nil {
		s.kafkaConsumer = consumer.NewConsumer(
			s.ConsumerGroup(),
			s.ConsumerGroupHandler(),
		)

		closer.Add(s.kafkaConsumer.Close)
	}

	return s.kafkaConsumer
}

// UserSaverConfig provides config for user saver service
func (s *ServiceProvider) UserSaverConfig() config.UserSaverConfig {
	if s.userSaverConfig == nil {
		conf, err := env.NewUserSaverConfigEnv()
		if err != nil {
			log.Fatalf("failed to load user saver config: %v", err)
		}

		s.userSaverConfig = conf
	}

	return s.userSaverConfig
}

// UserSaverService provices instance of user saver service
func (s *ServiceProvider) UserSaverService(ctx context.Context) service.ConsumerService {
	if s.userSaverService == nil {
		s.userSaverService = consumerService.NewUserSaverService(
			s.UserSaverConfig(),
			s.KafkaConsumer(),
			s.UserRepo(ctx),
			s.PasswordHasher(),
		)
	}

	return s.userSaverService
}

// UserServiceConfig provides config for user service
func (s *ServiceProvider) UserServiceConfig() config.UserServiceConfig {
	if s.userServiceConfig == nil {
		userServiceConfig, err := env.NewUserServiceConfigEnv()
		if err != nil {
			log.Fatalf("failed to load user service config: %v", err)
		}
		s.userServiceConfig = userServiceConfig
	}

	return s.userServiceConfig
}

// UserService initialises user service layer
func (s *ServiceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		var userCacheCli repository.UserCache
		if s.UserServiceConfig().UseCache() {
			userCacheCli = s.UserCache()
		}
		s.userService = userService.NewUserService(
			s.UserRepo(ctx),
			userCacheCli,
			s.TxManager(ctx),
			s.UserServiceConfig(),
		)
	}

	return s.userService
}

// AuthService initialises auth service layer
func (s *ServiceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewAuthService(
			s.UserRepo(ctx),
			s.RefreshTokenProvider(),
			s.AccessTokenProvider(),
			s.PasswordHasher(),
		)
	}

	return s.authService
}

// AccessService initialises access service layer
func (s *ServiceProvider) AccessService(ctx context.Context) service.AccessService {
	if s.accessService == nil {
		s.accessService = accessService.NewAccessService(s.AccessTokenProvider(), s.AccessRepo(ctx))
	}

	return s.accessService
}

// UserImpl Initialises user api server
func (s *ServiceProvider) UserImpl(ctx context.Context) *userImpl.Service {
	if s.userImpl == nil {
		s.userImpl = userImpl.NewService(s.UserService(ctx), s.PasswordHasher())
	}

	return s.userImpl
}

// AuthImpl Initialises auth api server
func (s *ServiceProvider) AuthImpl(ctx context.Context) *authImpl.Service {
	if s.authImpl == nil {
		s.authImpl = authImpl.NewService(s.AuthService(ctx))
	}

	return s.authImpl
}

// AccessImpl Initialises access api server
func (s *ServiceProvider) AccessImpl(ctx context.Context) *accessImpl.Service {
	if s.accessImpl == nil {
		s.accessImpl = accessImpl.NewService(s.AccessService(ctx))
	}

	return s.accessImpl
}
