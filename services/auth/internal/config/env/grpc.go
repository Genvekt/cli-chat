package env

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/Genvekt/cli-chat/services/auth/internal/config"
)

var _ config.GRPCConfig = (*gRPCConfigEnv)(nil)

const (
	grpcHostEnv = "GRPC_HOST"
	grpcPortEnv = "GRPC_PORT"

	grpcTLSEnabledEnv  = "GRPC_TLS_ENABLED"
	grpcTLSCertFileEnv = "GRPC_TLS_CERT_FILE"
	grpcTLSKeyFileEnv  = "GRPC_TLS_KEY_FILE"
)

type gRPCConfigEnv struct {
	host        string
	port        string
	tlsEnabled  bool
	tlsCertFile string
	tlsKeyFile  string
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

	isTLSEnabled, err := strconv.ParseBool(os.Getenv(grpcTLSEnabledEnv))
	if err != nil {
		return nil, fmt.Errorf("environment variable %q not set", grpcTLSEnabledEnv)
	}

	var tlsCertFile, tlsKeyFile string
	if isTLSEnabled {
		tlsCertFile = os.Getenv(grpcTLSCertFileEnv)
		if tlsCertFile == "" {
			return nil, fmt.Errorf("environment variable %q not set", grpcTLSCertFileEnv)
		}

		tlsKeyFile = os.Getenv(grpcTLSKeyFileEnv)
		if tlsKeyFile == "" {
			return nil, fmt.Errorf("environment variable %q not set", grpcTLSKeyFileEnv)
		}
	}

	return &gRPCConfigEnv{
		host:        host,
		port:        port,
		tlsEnabled:  isTLSEnabled,
		tlsCertFile: tlsCertFile,
		tlsKeyFile:  tlsKeyFile,
	}, nil
}

// Address provides host:port string
func (e *gRPCConfigEnv) Address() string {
	return net.JoinHostPort(e.host, e.port)
}

func (e *gRPCConfigEnv) TLSCertFile() string { return e.tlsCertFile }

func (e *gRPCConfigEnv) TLSKeyFile() string { return e.tlsKeyFile }

func (e *gRPCConfigEnv) IsTLSEnabled() bool { return e.tlsEnabled }
