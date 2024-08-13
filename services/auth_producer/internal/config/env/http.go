package env

import (
	"fmt"
	"net"
	"os"

	"github.com/Genvekt/cli-chat/services/auth_producer/internal/config"
)

var _ config.HTTPConfig = (*httpConfigEnv)(nil)

const (
	httpHostEnv = "HTTP_HOST"
	httpPortEnv = "HTTP_PORT"
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

// Address provides host:port string
func (e *httpConfigEnv) Address() string {
	return net.JoinHostPort(e.host, e.port)
}
