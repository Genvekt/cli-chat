package env

import (
	"fmt"
	"net"
	"os"

	config "github.com/Genvekt/cli-chat/services/chat-server/internal"
)

var _ config.GRPCConfig = (*gRPCConfigEnv)(nil)

const (
	grpcHostEnv = "GRPC_HOST"
	grpcPortEnv = "GRPC_PORT"
)

type gRPCConfigEnv struct {
	host string
	port string
}

// NewGRPCConfigEnv retrieves new gRPCConfigEnv instance
func NewGRPCConfigEnv() (*gRPCConfigEnv, error) {
	host := os.Getenv(grpcHostEnv)
	if host == "" {
		return nil, fmt.Errorf("environment variable %q not set", grpcHostEnv)
	}
	port := os.Getenv(grpcPortEnv)
	if port == "" {
		return nil, fmt.Errorf("environment variable %q not set", grpcPortEnv)
	}
	return &gRPCConfigEnv{host: host, port: port}, nil
}

// Address returns host:port string for grpc server
func (e *gRPCConfigEnv) Address() string {
	return net.JoinHostPort(e.host, e.port)
}
