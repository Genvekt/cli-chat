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

	config "github.com/Genvekt/cli-chat/services/chat-server/internal"
	"github.com/Genvekt/cli-chat/services/chat-server/internal/env"
	"github.com/Genvekt/cli-chat/services/chat-server/model"
	"github.com/Genvekt/cli-chat/services/chat-server/repository"
	"github.com/Genvekt/cli-chat/services/chat-server/repository/postgres"

	// use general reference for easier api version change in code
	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to .env config file")
}

type server struct {
	chatRepo repository.ChatRepository
	chatApi.UnimplementedChatV1Server
}

func (s *server) Create(_ context.Context, req *chatApi.CreateRequest) (*chatApi.CreateResponse, error) {
	log.Printf("Received create chat request: %+v", req)
	newChat := &model.Chat{}

	var err error
	newChat, err = s.chatRepo.Create(newChat)
	if err != nil {
		return nil, err
	}

	return &chatApi.CreateResponse{
		Id: newChat.Id,
	}, nil
}

func (s *server) SendMessage(_ context.Context, req *chatApi.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("Received send chat request: %+v", req)
	// Something here to send message...
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(_ context.Context, req *chatApi.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Received delete chat request: %+v", req)
	err := s.chatRepo.Delete(req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func runServer(grpcConf config.GRPCConfig, postgresConf config.PostgresConfig) error {
	ctx := context.Background()

	lis, err := net.Listen("tcp", grpcConf.Address())
	if err != nil {

		return fmt.Errorf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, postgresConf.DSN())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	chatRepo := postgres.NewChatPostgresRepository(ctx, pool)

	s := grpc.NewServer()
	reflection.Register(s)
	chatApi.RegisterChatV1Server(s, &server{chatRepo: chatRepo})

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
	err = runServer(grpcConfig, postgresConfig)
	if err != nil {
		log.Fatal(err)
	}
}
