package env

import (
	"fmt"
	"os"
	"strconv"
)

const (
	serviceNameEnv = "JAEGER_SERVICE_NAME"

	// use jaeger default envs
	samplerTypeEnv  = "JAEGER_SAMPLER_TYPE"
	samplerParamEnv = "JAEGER_SAMPLER_PARAM"
	agentAddressEnv = "JAEGER_AGENT_HOST"
)

type jaegerTracingConfigEnv struct {
	serviceName  string
	samplerType  string
	samplerParam float64
	agentAddress string
}

func NewJaegerTracingConfigEnv() (*jaegerTracingConfigEnv, error) {
	serviceName := os.Getenv(serviceNameEnv)
	if serviceName == "" {
		return nil, fmt.Errorf("environment variable %q not defined", serviceNameEnv)
	}

	samplerType := os.Getenv(samplerTypeEnv)
	if samplerType == "" {
		return nil, fmt.Errorf("environment variable %q not defined", samplerTypeEnv)
	}

	samplerParam, err := strconv.ParseFloat(os.Getenv(samplerParamEnv), 64)
	if err != nil {
		return nil, fmt.Errorf("environment variable %q not defined, float expected", samplerParamEnv)
	}

	agentAddress := os.Getenv(agentAddressEnv)
	if agentAddress == "" {
		return nil, fmt.Errorf("environment variable %q not defined", agentAddressEnv)
	}

	return &jaegerTracingConfigEnv{
		serviceName:  serviceName,
		samplerType:  samplerType,
		samplerParam: samplerParam,
		agentAddress: agentAddress,
	}, nil
}

func (j *jaegerTracingConfigEnv) ServiceName() string {
	return j.serviceName
}

func (j *jaegerTracingConfigEnv) SamplerType() string {
	return j.samplerType
}

func (j *jaegerTracingConfigEnv) SamplerParam() float64 {
	return j.samplerParam
}

func (j *jaegerTracingConfigEnv) AgentAddress() string {
	return j.agentAddress
}
