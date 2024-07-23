package config

import "github.com/joho/godotenv"

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

// PostgresConfig provides parameters related to Postgres database
type PostgresConfig interface {
	DSN() string
}
