package app

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Genvekt/cli-chat/libraries/db_client/pkg/db"
	"github.com/Genvekt/cli-chat/libraries/db_client/pkg/db/pg"
	"github.com/Genvekt/cli-chat/libraries/db_client/pkg/db/transaction"

	"github.com/Genvekt/cli-chat/libraries/closer/pkg/closer"

	chatImpl "github.com/Genvekt/cli-chat/services/chat-server/internal/api/chat"
	serviceClient "github.com/Genvekt/cli-chat/services/chat-server/internal/client/service"
	authClient "github.com/Genvekt/cli-chat/services/chat-server/internal/client/service/auth"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/config"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/config/env"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/repository"
	chatRepository "github.com/Genvekt/cli-chat/services/chat-server/internal/repository/chat"
	chatMemberRepository "github.com/Genvekt/cli-chat/services/chat-server/internal/repository/chat_member"
	messageRepository "github.com/Genvekt/cli-chat/services/chat-server/internal/repository/message"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/service"
	chatService "github.com/Genvekt/cli-chat/services/chat-server/internal/service/chat"
)

// ServiceProvider initialises and stores various dependencies as singletons
type ServiceProvider struct {
	gRPCConfig        config.GRPCConfig
	authCliGRPCConfig config.GRPCConfig
	postgresConfig    config.PostgresConfig

	dbClient  db.Client
	txManager db.TxManager

	chatRepo       repository.ChatRepository
	chatMemberRepo repository.ChatMemberRepository
	messageRepo    repository.MessageRepository

	authClient serviceClient.AuthClient

	chatService service.ChatService

	chatImpl *chatImpl.Service
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

// AuthCliGRPCConfig provides configuration parameters for auth service client
func (s *ServiceProvider) AuthCliGRPCConfig() config.GRPCConfig {
	if s.authCliGRPCConfig == nil {
		grpcConfig, err := env.NewAuthCliGRPCConfigEnv()
		if err != nil {
			log.Fatalf("failed to load auth grpc config: %v", err)
		}

		s.authCliGRPCConfig = grpcConfig
	}

	return s.authCliGRPCConfig
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
		dbClient, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to postgres: %v", err)
		}

		if err := dbClient.DB().Ping(ctx); err != nil {
			log.Fatalf("failed to ping postgres: %v", err)
		}

		closer.Add(func() error {
			return dbClient.Close()
		})

		s.dbClient = dbClient
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

// AuthClient provides auth service client dependency
func (s *ServiceProvider) AuthClient() serviceClient.AuthClient {
	if s.authClient == nil {
		var err error
		creds := insecure.NewCredentials()

		// configure TLS if it is enabled
		if s.AuthCliGRPCConfig().IsTLSEnabled() {
			creds, err = credentials.NewClientTLSFromFile(s.AuthCliGRPCConfig().TLSCertFile(), "")
			if err != nil {
				log.Fatalf("failed to load tls cert for auth grpc client: %v", err)
			}
			log.Println("Auth GRPC client: TLS enabled")
		}

		conn, err := grpc.NewClient(s.AuthCliGRPCConfig().Address(), grpc.WithTransportCredentials(creds))
		if err != nil {
			log.Fatalf("failed to connect to auth service: %v", err)
		}

		s.authClient = authClient.NewAuthClient(authClient.NewAuthGrpcClient(conn))

		closer.Add(conn.Close)
	}

	return s.authClient
}

// ChatRepo provides chat repository dependency
func (s *ServiceProvider) ChatRepo(ctx context.Context) repository.ChatRepository {
	if s.chatRepo == nil {
		s.chatRepo = chatRepository.NewChatPostgresRepository(s.DBClient(ctx))
	}

	return s.chatRepo
}

// ChatMemberRepo provides chat member repository dependency
func (s *ServiceProvider) ChatMemberRepo(ctx context.Context) repository.ChatMemberRepository {
	if s.chatMemberRepo == nil {
		s.chatMemberRepo = chatMemberRepository.NewChatMemberPostgresRepository(s.DBClient(ctx))
	}

	return s.chatMemberRepo
}

// MessageRepo provides message repository dependency
func (s *ServiceProvider) MessageRepo(ctx context.Context) repository.MessageRepository {
	if s.messageRepo == nil {
		s.messageRepo = messageRepository.NewMessagePostgresRepository(s.DBClient(ctx))
	}

	return s.messageRepo
}

// ChatService initialises chat service layer
func (s *ServiceProvider) ChatService(ctx context.Context) service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewChatService(
			s.ChatRepo(ctx),
			s.ChatMemberRepo(ctx),
			s.MessageRepo(ctx),
			s.AuthClient(),
			s.TxManager(ctx),
		)
	}

	return s.chatService
}

// ChatImpl Initialises chat api server
func (s *ServiceProvider) ChatImpl(ctx context.Context) *chatImpl.Service {
	if s.chatImpl == nil {
		s.chatImpl = chatImpl.NewService(s.ChatService(ctx))
	}

	return s.chatImpl
}
