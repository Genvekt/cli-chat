package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	// use general reference for easier api version change in code
	userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
	config "github.com/Genvekt/cli-chat/services/auth/internal"
	"github.com/Genvekt/cli-chat/services/auth/internal/env"
	"github.com/Genvekt/cli-chat/services/auth/model"
	"github.com/Genvekt/cli-chat/services/auth/repository"
	"github.com/Genvekt/cli-chat/services/auth/repository/postgres"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to .env config file")
}

type server struct {
	userRepo repository.UserRepository
	userApi.UnimplementedUserV1Server
}

func (s *server) Get(ctx context.Context, req *userApi.GetRequest) (*userApi.GetResponse, error) {
	dbUser, err := s.userRepo.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &userApi.GetResponse{
		Id:        dbUser.ID,
		Name:      dbUser.Name,
		Email:     dbUser.Email,
		Role:      userApi.UserRole(dbUser.Role),
		CreatedAt: timestamppb.New(dbUser.CreatedAt),
		UpdatedAt: timestamppb.New(dbUser.UpdatedAt),
	}, nil
}

func (s *server) Create(ctx context.Context, req *userApi.CreateRequest) (*userApi.CreateResponse, error) {
	newUser := &model.User{
		Name:  req.Name,
		Email: req.Email,
		Role:  int(req.Role),
	}

	err := s.userRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &userApi.CreateResponse{
		Id: newUser.ID,
	}, nil
}

func (s *server) Update(ctx context.Context, req *userApi.UpdateRequest) (*emptypb.Empty, error) {
	updateFunc := func(user *model.User) error {
		if req.Email != nil {
			user.Email = req.Email.Value
		}
		if req.Name != nil {
			user.Name = req.Name.Value
		}
		if req.Role != nil {
			user.Role = int(*req.Role)
		}
		return nil
	}

	err := s.userRepo.Update(ctx, req.Id, updateFunc)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *userApi.DeleteRequest) (*emptypb.Empty, error) {
	err := s.userRepo.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func runServer(ctx context.Context, grpcConf config.GRPCConfig, postgresConf config.PostgresConfig) error {
	lis, err := net.Listen("tcp", grpcConf.Address())
	if err != nil {

		return fmt.Errorf("failed to listen: %v", err)
	}

	pool, err := pgxpool.Connect(ctx, postgresConf.DSN())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	userRepo := postgres.NewUserRepositoryPostgres(pool)

	s := grpc.NewServer()
	reflection.Register(s)
	userApi.RegisterUserV1Server(s, &server{userRepo: userRepo})

	log.Printf("started gRPC server at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

func main() {
	flag.Parse()
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfigEnv()
	if err != nil {
		log.Fatalf("failed to load grpc config: %v", err)
	}

	postgresConfig, err := env.NewPostgresConfigEnv()
	if err != nil {
		log.Fatalf("failed to load postgres config: %v", err)
	}

	ctx := context.Background()

	err = runServer(ctx, grpcConfig, postgresConfig)
	if err != nil {
		log.Fatal(err)
	}
}
