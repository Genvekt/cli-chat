package settings

import (
	"os"
	"strconv"
)

const (
	grpcDefaultPort = 50051
	grpcPostEnv     = "GRPC_PORT"
)

var settingsObj *settings

// settings struct holds all env parameters
type settings struct {
	GrpcPort int
}

// setDefaultValues applies default values to unfilled envs
func (s *settings) setDefaultValues() {
	_, isGrpcPortSet := os.LookupEnv(grpcPostEnv)
	if !isGrpcPortSet {
		s.GrpcPort = grpcDefaultPort
	}
}

// GetSettings crates settings singleton and used to retrieve it
func GetSettings() *settings {
	if settingsObj == nil {
		grpcPort, _ := strconv.Atoi(os.Getenv(grpcPostEnv))
		settingsObj = &settings{
			GrpcPort: grpcPort,
		}
		settingsObj.setDefaultValues()
	}

	return settingsObj
}
