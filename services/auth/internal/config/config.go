package config

import (
	"time"

	"github.com/joho/godotenv"
)

// Load reads .env file into environment vars
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

// GRPCConfig provides parameters related to GRPC server
type GRPCConfig interface {
	Address() string
}

// HTTPConfig provides parameters related to HTTP server
type HTTPConfig interface {
	Address() string
}

// PostgresConfig provides parameters related to Postgres database
type PostgresConfig interface {
	DSN() string
}

// UserServiceConfig provides parameters related to user service
type UserServiceConfig interface {
	CacheTTL() time.Duration
	NoCache() bool
	UseCache() bool
}
