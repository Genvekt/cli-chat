package app

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Genvekt/cli-chat/libraries/closer/pkg/closer"
	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/cli"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/config/ini"

	chatCli "github.com/Genvekt/cli-chat/services/chat-client/internal/cli"
	serviceClient "github.com/Genvekt/cli-chat/services/chat-client/internal/client/service"
	authClient "github.com/Genvekt/cli-chat/services/chat-client/internal/client/service/auth"
	chatClient "github.com/Genvekt/cli-chat/services/chat-client/internal/client/service/chat"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/config"
	"github.com/Genvekt/cli-chat/services/chat-client/internal/config/env"
)

type ServiceProvider struct {
	iniConfPath string
	profileName string

	chatConnection *grpc.ClientConn
	chatGRPCConfig config.GRPCConfig
	chatClient     serviceClient.ChatClient

	authConnection *grpc.ClientConn
	authGRPCConfig config.GRPCConfig
	authClient     serviceClient.AuthClient

	profileConfig  config.ProfileConfig
	chatCliService *cli.CliService
}

func newServiceProvider(iniConfPath string, profileName string) *ServiceProvider {
	return &ServiceProvider{
		iniConfPath: iniConfPath,
		profileName: profileName,
	}
}

func (s *ServiceProvider) ChatGRPCConfig() config.GRPCConfig {
	if s.chatGRPCConfig == nil {
		conf, err := env.NewChatGRPCConfigEnv()
		if err != nil {
			logger.Fatal(
				"cannot load chat client env vars: %v",
				zap.Error(err),
			)
		}

		s.chatGRPCConfig = conf
	}

	return s.chatGRPCConfig
}

func (s *ServiceProvider) AuthGRPCConfig() config.GRPCConfig {
	if s.authGRPCConfig == nil {
		conf, err := env.NewAuthGRPCConfigEnv()
		if err != nil {
			logger.Fatal(
				"cannot load chat client env vars: %v",
				zap.Error(err),
			)
		}

		s.authGRPCConfig = conf
	}

	return s.authGRPCConfig
}

// ChatConnection provides grpc connection to chat service
func (s *ServiceProvider) ChatConnection() *grpc.ClientConn {
	if s.chatConnection == nil {
		var err error
		creds := insecure.NewCredentials()

		conn, err := grpc.NewClient(
			s.ChatGRPCConfig().Address(),
			grpc.WithTransportCredentials(creds),
		)
		if err != nil {
			logger.Fatal("failed to connect to chat service", zap.Error(err))
		}

		closer.Add(conn.Close)

		s.chatConnection = conn
	}

	return s.chatConnection
}

// AuthConnection provides grpc connection to auth service
func (s *ServiceProvider) AuthConnection() *grpc.ClientConn {
	if s.authConnection == nil {
		var err error
		creds := insecure.NewCredentials()

		conn, err := grpc.NewClient(
			s.AuthGRPCConfig().Address(),
			grpc.WithTransportCredentials(creds),
		)
		if err != nil {
			logger.Fatal("failed to connect to auth service", zap.Error(err))
		}

		closer.Add(conn.Close)

		s.authConnection = conn
	}

	return s.authConnection
}

// ChatClient provides chat service client dependency
func (s *ServiceProvider) ChatClient() serviceClient.ChatClient {
	if s.chatClient == nil {
		s.chatClient = chatClient.NewChatClient(chatClient.NewChatGrpcClientWrapper(s.ChatConnection()))
	}

	return s.chatClient
}

// AuthClient provides auth service client dependency
func (s *ServiceProvider) AuthClient() serviceClient.AuthClient {
	if s.authClient == nil {
		s.authClient = authClient.NewAuthClient(authClient.NewAuthGrpcClientWrapper(s.AuthConnection()))
	}

	return s.authClient
}

func (s *ServiceProvider) ProfileConfig(configPath string, profileName string) config.ProfileConfig {
	if s.profileConfig == nil {
		if profileName == "" {
			// try to get profile from env before referencing ini config
			confEnv, err := env.NewProfileConfigEnv()
			if err == nil {
				s.profileConfig = confEnv
			}
		}

		if s.profileConfig == nil {
			// try to get profile from ini config
			confIni, err := ini.NewProfileConfigIni(configPath, profileName)
			if err != nil {
				logger.Fatal("cannot load profile data", zap.Error(err))
			}

			s.profileConfig = confIni

		}
	}

	return s.profileConfig
}

// ChatCliService provides chat cli application dependency
func (s *ServiceProvider) ChatCliService() *cli.CliService {
	if s.chatCliService == nil {
		s.chatCliService = chatCli.NewChatCliService(
			s.ProfileConfig(s.iniConfPath, s.profileName),
			s.ChatClient(),
			s.AuthClient(),
		)
	}

	return s.chatCliService
}
