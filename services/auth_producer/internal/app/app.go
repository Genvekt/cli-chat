package app

import (
	"context"
	"flag"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Genvekt/cli-chat/libraries/closer/pkg/closer"
	"github.com/Genvekt/cli-chat/services/auth_producer/internal/config"
)

// App is an application starting point
type App struct {
	configPath string
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
	log.Printf("Started HTTP server at %v", a.provider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
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
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := a.runHTTPServer(); err != nil {
			log.Fatalf("Failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.runProducers(ctx); err != nil {
			log.Fatalf("Failed to run consumers: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) runProducers(ctx context.Context) error {

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := a.provider.UserCreatorService().RunProducer(ctx)
		if err != nil {
			log.Printf("failed to run procuser: %s", err.Error())
		}
	}()

	wg.Wait()

	return nil
}
