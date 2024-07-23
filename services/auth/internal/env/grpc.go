package env

import (
	"fmt"
	"net"
	"os"

	config "github.com/Genvekt/cli-chat/services/auth/internal"
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

func (e *gRPCConfigEnv) Address() string {
	return net.JoinHostPort(e.host, e.port)
}
