package config

import "github.com/joho/godotenv"

// LoadEnv reads .env file into environment vars
func LoadEnv(filePath string) error {
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

// GRPCConfig provides envs related to grpc client
type GRPCConfig interface {
	Address() string
}

type ProfileConfig interface {
	Username() string
	Password() string
}
