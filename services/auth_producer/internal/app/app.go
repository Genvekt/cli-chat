package app

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Genvekt/cli-chat/libraries/closer/pkg/closer"
	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
	"github.com/Genvekt/cli-chat/services/auth_producer/internal/config"
)

// App is an application starting point
type App struct {
	configPath string
	logLevel   string
	provider   *ServiceProvider
	httpServer *http.Server
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
		a.initServiceProvider,
		a.initHTTPServer,
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

func (a *App) initHTTPServer(_ context.Context) error {
	http.HandleFunc("/", a.provider.UserAPI().HandleCreate)

	a.httpServer = &http.Server{
		Addr:              a.provider.HTTPConfig().Address(),
		Handler:           http.DefaultServeMux,
		ReadHeaderTimeout: time.Second * 5,
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

// Run starts application and triggers closer on stop
func (a *App) Run(_ context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := a.runHTTPServer(); err != nil {
			logger.Fatal("Failed to run HTTP server", zap.Error(err))
		}
	}()

	wg.Wait()

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
