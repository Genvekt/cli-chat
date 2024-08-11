package env

import (
	"fmt"
	"net"
	"os"

	"github.com/Genvekt/cli-chat/services/auth/internal/config"
)

var _ config.GRPCConfig = (*gRPCConfigEnv)(nil)

const (
	httpHostEnv = "HTTP_HOST"
	httpPortEnv = "HTTP_PORT"

	swaggerHostEnv = "SWAGGER_HOST"
	swaggerPortEnv = "SWAGGER_PORT"
)

type httpConfigEnv struct {
	host string
	port string
}

// NewHTTPConfigEnv retrieves new httpConfigEnv instance
func NewHTTPConfigEnv() (*httpConfigEnv, error) {
	host := os.Getenv(httpHostEnv)
	if host == "" {
		return nil, fmt.Errorf("environment variable %q not set", httpHostEnv)
	}

	port := os.Getenv(httpPortEnv)
	if port == "" {
		return nil, fmt.Errorf("environment variable %q not set", httpPortEnv)
	}

	return &httpConfigEnv{host: host, port: port}, nil
}

func NewSwaggerConfigEnv() (*httpConfigEnv, error) {
	host := os.Getenv(swaggerHostEnv)
	if host == "" {
		return nil, fmt.Errorf("environment variable %q not set", swaggerHostEnv)
	}

	port := os.Getenv(swaggerPortEnv)
	if port == "" {
		return nil, fmt.Errorf("environment variable %q not set", swaggerPortEnv)
	}

	return &httpConfigEnv{host: host, port: port}, nil
}

// Address provides host:port string
func (e *httpConfigEnv) Address() string {
	return net.JoinHostPort(e.host, e.port)
}
