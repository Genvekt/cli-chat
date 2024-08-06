package app

import (
	"context"
	"log"

	userImpl "github.com/Genvekt/cli-chat/services/auth/internal/api/user"
	"github.com/Genvekt/cli-chat/services/auth/internal/client/db"
	"github.com/Genvekt/cli-chat/services/auth/internal/client/db/pg"
	"github.com/Genvekt/cli-chat/services/auth/internal/client/db/transaction"
	"github.com/Genvekt/cli-chat/services/auth/internal/closer"
	"github.com/Genvekt/cli-chat/services/auth/internal/config"
	"github.com/Genvekt/cli-chat/services/auth/internal/config/env"
	"github.com/Genvekt/cli-chat/services/auth/internal/repository"
	userRepository "github.com/Genvekt/cli-chat/services/auth/internal/repository/user"
	"github.com/Genvekt/cli-chat/services/auth/internal/service"
	userService "github.com/Genvekt/cli-chat/services/auth/internal/service/user"
)

// ServiceProvider initialises and stores various dependencies as singletons
type ServiceProvider struct {
	gRPCConfig     config.GRPCConfig
	postgresConfig config.PostgresConfig

	dbClient  db.Client
	txManager db.TxManager

	userRepo repository.UserRepository

	userService service.UserService

	userImpl *userImpl.Service
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

// UserRepo provides user repository dependency
func (s *ServiceProvider) UserRepo(ctx context.Context) repository.UserRepository {
	if s.userRepo == nil {
		s.userRepo = userRepository.NewUserRepositoryPostgres(s.DBClient(ctx))
	}

	return s.userRepo
}

// UserService initialises user service layer
func (s *ServiceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewUserService(s.UserRepo(ctx), s.TxManager(ctx))
	}

	return s.userService
}

// UserImpl Initialises user api server
func (s *ServiceProvider) UserImpl(ctx context.Context) *userImpl.Service {
	if s.userImpl == nil {
		s.userImpl = userImpl.NewService(s.UserService(ctx))
	}

	return s.userImpl
}
