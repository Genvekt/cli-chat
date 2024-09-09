package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/credentials"

	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
	"github.com/Genvekt/cli-chat/services/auth/internal/metric"

	"github.com/rakyll/statik/fs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	accessApi "github.com/Genvekt/cli-chat/libraries/api/access/v1"
	authApi "github.com/Genvekt/cli-chat/libraries/api/auth/v1"
	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/libraries/closer/pkg/closer"
	"github.com/Genvekt/cli-chat/services/auth/internal/config"
	"github.com/Genvekt/cli-chat/services/auth/internal/interceptor"
)

// App is an application starting point
type App struct {
	configPath       string
	logLevel         string
	provider         *ServiceProvider
	grpcServer       *grpc.Server
	httpServer       *http.Server
	swaggerServer    *http.Server
	prometheusServer *http.Server
}

// NewApp initialises app and all its dependencies
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
		a.initMetrics,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
		a.initPrometheusServer,
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

func (a *App) initMetrics(ctx context.Context) error {
	err := metric.Init(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.provider = newServiceProvider()

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	var err error

	// configure TLS if it is enabled
	creds := insecure.NewCredentials()
	if a.provider.GRPCConfig().IsTLSEnabled() {
		creds, err = credentials.NewServerTLSFromFile(a.provider.GRPCConfig().TLSCertFile(), a.provider.GRPCConfig().TLSKeyFile())
		if err != nil {
			return err
		}
		logger.Debug("GRPC TLS enabled")
	}

	a.grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.ChainUnaryInterceptor(
			interceptor.LogInterceptor,
			interceptor.MetricsInterceptor,
			interceptor.ValidateInterceptor,
		),
	)

	reflection.Register(a.grpcServer)

	userApi.RegisterUserV1Server(a.grpcServer, a.provider.UserImpl(ctx))
	authApi.RegisterAuthV1Server(a.grpcServer, a.provider.AuthImpl(ctx))
	accessApi.RegisterAccessV1Server(a.grpcServer, a.provider.AccessImpl(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := userApi.RegisterUserV1HandlerFromEndpoint(ctx, mux, a.provider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:              a.provider.HTTPConfig().Address(),
		Handler:           corsMiddleware.Handler(mux),
		ReadHeaderTimeout: time.Second * 5,
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/userApi.swagger.json", serveSwaggerFile("/userApi.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:              a.provider.SwaggerConfig().Address(),
		Handler:           mux,
		ReadHeaderTimeout: time.Second * 5,
	}

	return nil
}

func (a *App) initPrometheusServer(_ context.Context) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr:              a.provider.PrometheusConfig().Address(),
		Handler:           mux,
		ReadHeaderTimeout: time.Second * 5,
	}

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

func (a *App) runHTTPServer() error {
	logger.Info("Started HTTP server",
		zap.String("address", a.provider.HTTPConfig().Address()),
	)

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	logger.Info("Started Swagger server",
		zap.String("address", a.provider.SwaggerConfig().Address()),
	)

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runPrometheusServer() error {
	logger.Info("Started Prometheus server",
		zap.String("address", a.provider.PrometheusConfig().Address()),
	)

	err := a.prometheusServer.ListenAndServe()
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

	wg := sync.WaitGroup{}
	wg.Add(5)

	go func() {
		defer wg.Done()
		if err := a.runGRPCServer(ctx); err != nil {
			logger.Fatal("Failed to run gRPC server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.runHTTPServer(); err != nil {
			logger.Fatal("Failed to run HTTP server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.runSwaggerServer(); err != nil {
			logger.Fatal("Failed to run swagger server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.runPrometheusServer(); err != nil {
			logger.Fatal("Failed to run prometheus server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.runConsumers(ctx); err != nil {
			logger.Fatal("Failed to run consumers", zap.Error(err))
		}
	}()

	wg.Wait()

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		logger.Debug("Serving swagger file", zap.String("path", path))

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Debug("Open swagger file", zap.String("path", path))

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			closeErr := file.Close()
			if closeErr != nil {
				logger.Error("Failed to close file",
					zap.Error(closeErr),
					zap.String("path", path),
				)
			}
		}()

		logger.Debug("Read swagger file", zap.String("path", path))

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Debug("Write swagger file", zap.String("path", path))

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Debug("Served swagger file", zap.String("path", path))
	}
}

func (a *App) runConsumers(ctx context.Context) error {

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		logger.Info("Started user saver consumer")

		err := a.provider.UserSaverService(ctx).RunConsumer(ctx)
		if err != nil {
			logger.Error("Failed to run consumer", zap.Error(err))
		}
	}()

	return nil
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
