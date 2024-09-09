package app

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/Genvekt/cli-chat/libraries/closer/pkg/closer"
	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/interceptor"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/tracing"

	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/config"
)

// App is an entrypoint of this application
type App struct {
	logLevel   string
	configPath string
	provider   *ServiceProvider
	grpcServer *grpc.Server
}

// NewApp initialises app with all dependencies
func NewApp(ctx context.Context) (*App, error) {
	application := &App{}

	err := application.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return application, nil
}

func (a *App) initDeps(ctx context.Context) error {
	deps := []func(context.Context) error{
		a.initArgs,
		a.initConfig,
		a.initLogger,
		a.initServiceProvider,
		a.initJaegerTracing,
		a.initGRPCServer,
	}

	for _, dep := range deps {
		if err := dep(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initArgs(_ context.Context) error {
	flag.StringVar(&a.configPath, "config-path", ".env", "path to .env config file")
	flag.StringVar(&a.logLevel, "l", "info", "log level")
	flag.Parse()

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(a.configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logLevel, err := a.getLogAtomicLevel()
	if err != nil {
		return err
	}
	logger.Init(a.getLogCore(*logLevel))
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.provider = newServiceProvider()

	return nil
}

func (a *App) initJaegerTracing(_ context.Context) error {
	err := tracing.Init(a.provider.JaegerConfig())
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	authorization := interceptor.NewAuthorization(a.provider.AccessClient())
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			interceptor.ServerTracingInterceptor,
			interceptor.ValidateInterceptor,
			authorization.Interceptor,
		),
	)

	reflection.Register(a.grpcServer)

	chatApi.RegisterChatV1Server(a.grpcServer, a.provider.ChatImpl(ctx))

	return nil
}

func (a *App) runGRPCServer(_ context.Context) error {
	lis, err := net.Listen("tcp", a.provider.GRPCConfig().Address())
	if err != nil {

		return fmt.Errorf("failed to listen: %v", err)
	}

	logger.Info("Started gRPC server",
		zap.String("address", lis.Addr().String()),
	)

	err = a.grpcServer.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

// Run starts application and triggers closer on stop
func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runGRPCServer(ctx)
}

func (a *App) getLogCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func (a *App) getLogAtomicLevel() (*zap.AtomicLevel, error) {
	var level zapcore.Level
	if err := level.Set(a.logLevel); err != nil {
		return nil, fmt.Errorf("failed to set log level: %v", err)
	}
	atomicLevel := zap.NewAtomicLevelAt(level)

	return &atomicLevel, nil
}
