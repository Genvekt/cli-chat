package config

import "github.com/joho/godotenv"

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

type GRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	DSN() string
}
