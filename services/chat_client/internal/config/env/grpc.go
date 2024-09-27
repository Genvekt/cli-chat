package env

import (
	"fmt"
	"net"
	"os"

	"github.com/Genvekt/cli-chat/services/chat-client/internal/config"
)

var _ config.GRPCConfig = (*gRPCConfigEnv)(nil)

const (
	chatGrpcHostEnv = "CHAT_GRPC_HOST"
	chatGrpcPortEnv = "CHAT_GRPC_PORT"

	authGrpcHostEnv = "AUTH_GRPC_HOST"
	authGrpcPortEnv = "AUTH_GRPC_PORT"
)

type gRPCConfigEnv struct {
	host        string
	port        string
	tlsEnabled  bool
	tlsCertFile string
	tlsKeyFile  string
}

// NewChatGRPCConfigEnv retrieves grpc config for chat service connection
func NewChatGRPCConfigEnv() (*gRPCConfigEnv, error) {
	host := os.Getenv(chatGrpcHostEnv)
	if host == "" {
		return nil, fmt.Errorf("environment variable %q not set", chatGrpcHostEnv)
	}

	port := os.Getenv(chatGrpcPortEnv)
	if port == "" {
		return nil, fmt.Errorf("environment variable %q not set", chatGrpcPortEnv)
	}

	return &gRPCConfigEnv{host: host, port: port}, nil
}

// NewAuthGRPCConfigEnv retrieves grpc config for auth service connection
func NewAuthGRPCConfigEnv() (*gRPCConfigEnv, error) {
	host := os.Getenv(authGrpcHostEnv)
	if host == "" {
		return nil, fmt.Errorf("environment variable %q not set", authGrpcHostEnv)
	}

	port := os.Getenv(authGrpcPortEnv)
	if port == "" {
		return nil, fmt.Errorf("environment variable %q not set", authGrpcPortEnv)
	}

	return &gRPCConfigEnv{host: host, port: port}, nil
}

// Address returns host:port string for grpc server
func (e *gRPCConfigEnv) Address() string {
	return net.JoinHostPort(e.host, e.port)
}
