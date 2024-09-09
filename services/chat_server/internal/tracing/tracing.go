package tracing

import (
	jaegerCfg "github.com/uber/jaeger-client-go/config"

	"github.com/Genvekt/cli-chat/services/chat-server/internal/config"
)

func Init(cfg config.JaegerTracingConfig) error {
	tracingCfg := jaegerCfg.Configuration{
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  cfg.SamplerType(),
			Param: cfg.SamplerParam(),
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LocalAgentHostPort: cfg.AgentAddress(),
		},
	}

	_, err := tracingCfg.InitGlobalTracer(cfg.ServiceName())
	if err != nil {
		return err
	}

	return nil
}
