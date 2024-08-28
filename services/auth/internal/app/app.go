package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/rs/cors"
	"google.golang.org/grpc/credentials"

	"github.com/rakyll/statik/fs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	"github.com/Genvekt/cli-chat/libraries/closer/pkg/closer"
	"github.com/Genvekt/cli-chat/services/auth/internal/config"
	"github.com/Genvekt/cli-chat/services/auth/internal/interceptor"
)

// App is an application starting point
type App struct {
	configPath    string
	provider      *ServiceProvider
	grpcServer    *grpc.Server
	httpServer    *http.Server
	swaggerServer *http.Server
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
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
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
	var err error

	// configure TLS if it is enabled
	creds := insecure.NewCredentials()
	if a.provider.GRPCConfig().IsTLSEnabled() {
		creds, err = credentials.NewServerTLSFromFile(a.provider.GRPCConfig().TLSCertFile(), a.provider.GRPCConfig().TLSKeyFile())
		if err != nil {
			return err
		}
		log.Println("GRPC TLS enabled")
	}

	a.grpcServer = grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

	reflection.Register(a.grpcServer)

	userApi.RegisterUserV1Server(a.grpcServer, a.provider.UserImpl(ctx))

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

func (a *App) runHTTPServer() error {
	log.Printf("Started HTTP server at %v", a.provider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Started Swagger server at %s", a.provider.SwaggerConfig().Address())

	err := a.swaggerServer.ListenAndServe()
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
	wg.Add(4)

	go func() {
		defer wg.Done()
		if err := a.runGRPCServer(ctx); err != nil {
			log.Fatalf("Failed to run gRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.runHTTPServer(); err != nil {
			log.Fatalf("Failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.runSwaggerServer(); err != nil {
			log.Fatalf("Failed to run swagger server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := a.runConsumers(ctx); err != nil {
			log.Fatalf("Failed to run consumers: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			closeErr := file.Close()
			if closeErr != nil {
				log.Printf("Failed to close file: %v", closeErr)
			}
		}()

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}

func (a *App) runConsumers(ctx context.Context) error {

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		log.Printf("Started user saver consumer")

		err := a.provider.UserSaverService(ctx).RunConsumer(ctx)
		if err != nil {
			log.Printf("failed to run consumer: %s", err.Error())
		}
	}()

	return nil
}
