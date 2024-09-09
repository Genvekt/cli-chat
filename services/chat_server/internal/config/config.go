package config

import "github.com/joho/godotenv"

// Load reads .env file into environment variables
func Load(filePath string) error {
	if filePath == "" {
		// nothing to load
		return nil
	}

	err := godotenv.Load(filePath)
	if err != nil {
		return err
	}

	return nil
}

// GRPCConfig provides envs related to grpc server
type GRPCConfig interface {
	Address() string
	TLSCertFile() string
	TLSKeyFile() string
	IsTLSEnabled() bool
}

// PostgresConfig provides envs related to postgres db
type PostgresConfig interface {
	DSN() string
}

// JaegerTracingConfig provides envs related to jaeger
type JaegerTracingConfig interface {
	ServiceName() string
	SamplerType() string
	SamplerParam() float64
	AgentAddress() string
}
