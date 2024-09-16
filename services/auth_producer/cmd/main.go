package main

import (
	"context"

	"go.uber.org/zap"

	"github.com/Genvekt/cli-chat/libraries/logger/pkg/logger"
	"github.com/Genvekt/cli-chat/services/auth_producer/internal/app"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx)
	if err != nil {
		logger.Fatal("failed to initialize application", zap.Error(err))
	}

	err = application.Run(ctx)
	if err != nil {
		logger.Fatal("failed to run application", zap.Error(err))
	}
}
