package env

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/config"
)

var _ config.GRPCConfig = (*gRPCConfigEnv)(nil)

const (
	grpcHostEnv = "GRPC_HOST"
	grpcPortEnv = "GRPC_PORT"

	authGrpcHostEnv        = "AUTH_GRPC_HOST"
	authGrpcPortEnv        = "AUTH_GRPC_PORT"
	authGrpcTLSEnabledEnv  = "AUTH_GRPC_TLS_ENABLED"
	authGrpcTLSCertFileEnv = "AUTH_GRPC_TLS_CERT_FILE"
)

type gRPCConfigEnv struct {
	host        string
	port        string
	tlsEnabled  bool
	tlsCertFile string
	tlsKeyFile  string
}

// NewAuthCliGRPCConfigEnv retrieves grpc config for auth service connection
func NewAuthCliGRPCConfigEnv() (*gRPCConfigEnv, error) {
	host := os.Getenv(authGrpcHostEnv)
	if host == "" {
		return nil, fmt.Errorf("environment variable %q not set", authGrpcHostEnv)
	}

	port := os.Getenv(authGrpcPortEnv)
	if port == "" {
		return nil, fmt.Errorf("environment variable %q not set", authGrpcPortEnv)
	}

	isTLSEnabled, err := strconv.ParseBool(os.Getenv(authGrpcTLSEnabledEnv))
	if err != nil {
		return nil, fmt.Errorf("environment variable %q not set", authGrpcTLSEnabledEnv)
	}

	var tlsCertFile string
	if isTLSEnabled {
		tlsCertFile = os.Getenv(authGrpcTLSCertFileEnv)
		if tlsCertFile == "" {
			return nil, fmt.Errorf("environment variable %q not set", authGrpcTLSCertFileEnv)
		}
	}

	return &gRPCConfigEnv{host: host, port: port, tlsEnabled: isTLSEnabled, tlsCertFile: tlsCertFile}, nil
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

	return &gRPCConfigEnv{host: host, port: port, tlsEnabled: false}, nil
}

// Address returns host:port string for grpc server
func (e *gRPCConfigEnv) Address() string {
	return net.JoinHostPort(e.host, e.port)
}

func (e *gRPCConfigEnv) TLSCertFile() string { return e.tlsCertFile }

func (e *gRPCConfigEnv) TLSKeyFile() string { return e.tlsKeyFile }

func (e *gRPCConfigEnv) IsTLSEnabled() bool { return e.tlsEnabled }
