package main

import (
	"context"
	"log"

	"github.com/Genvekt/cli-chat/services/auth/internal/app"

	// make init from static swagger service visible
	_ "github.com/Genvekt/cli-chat/services/auth/statik"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to initialize application: %v", err)
	}

	err = application.Run(ctx)
	if err != nil {
		log.Fatalf("failed to run application: %v", err)
	}
}
