package app

import (
  "context"
  "fmt"
  "os"

  "github.com/natefinch/lumberjack"
  "github.com/spf13/cobra"
  "go.uber.org/zap"
  "go.uber.org/zap/zapcore"

  "github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
  "github.com/Genvekt/cli-chat/services/chat-client/internal/cli"
  "github.com/Genvekt/cli-chat/services/chat-client/internal/config"
)

var (
  profileName string
  iniConfig   string
  envConfig   string

  application *App
)

type App struct {
  provider   *ServiceProvider
  cliRoot    *cobra.Command
  cliService cli.CliService
}

// InitApp initialises app and all its dependencies
func InitApp(ctx context.Context, profile string, iniConf string, envConf string) error {
  profileName = profile
  iniConfig = iniConf
  envConfig = envConf

  application = &App{}

  err := application.initDeps(ctx)
  if err != nil {
    return err
  }

  return nil
}

func (a *App) initDeps(ctx context.Context) error {
  deps := []func(context.Context) error{
    a.initConfig,
    a.initLogger,
    a.initServiceProvider,
  }

  for _, dep := range deps {
    if err := dep(ctx); err != nil {
      return err
    }
  }

  return nil
}

func (a *App) initConfig(_ context.Context) error {
  err := config.LoadEnv(envConfig)
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
  a.provider = newServiceProvider(iniConfig, profileName)

  return nil
}

func CreateChat(ctx context.Context, name string, usernames []string) error {
  return application.provider.ChatCliService().CreateChat(ctx, name, usernames)
}

func Connect(ctx context.Context, chatID int64) error {
  return application.provider.ChatCliService().Connect(ctx, chatID)
}

func SendMessage(ctx context.Context, chatID int64, message string) error {
  return application.provider.ChatCliService().SendMessage(ctx, chatID, message)
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
  if err := level.Set(zapcore.InfoLevel.String()); err != nil {
    return nil, fmt.Errorf("failed to set log level: %v", err)
  }
  atomicLevel := zap.NewAtomicLevelAt(level)

  return &atomicLevel, nil
}
