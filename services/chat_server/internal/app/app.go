package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/Genvekt/cli-chat/libraries/closer/pkg/closer"

	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/config"
)

// App is an entrypoint of this application
type App struct {
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
		a.initServiceProvider,
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

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	reflection.Register(a.grpcServer)
	chatApi.RegisterChatV1Server(a.grpcServer, a.provider.ChatImpl(ctx))

	return nil
}

func (a *App) runGRPCServer(_ context.Context) error {
	lis, err := net.Listen("tcp", a.provider.GRPCConfig().Address())
	if err != nil {

		return fmt.Errorf("failed to listen: %v", err)
	}

	log.Printf("Started gRPC server at %v", lis.Addr())

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
