package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Genvekt/cli-chat/services/chat-server/settings"

	// use general reference for easier api version change in code
	chatApi "github.com/Genvekt/cli-chat/libraries/api/chat/v1"
)

type server struct {
	chatApi.UnimplementedChatV1Server
}

func (s *server) Create(_ context.Context, req *chatApi.CreateRequest) (*chatApi.CreateResponse, error) {
	log.Printf("Received create chat request: %+v", req)
	return &chatApi.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) SendMessage(_ context.Context, req *chatApi.SendMessageRequest) (*emptypb.Empty, error) {
	log.Printf("Received send chat request: %+v", req)
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(_ context.Context, req *chatApi.DeleteRequest) (*emptypb.Empty, error) {
	log.Printf("Received delete chat request: %+v", req)
	return &emptypb.Empty{}, nil
}

func main() {
	env := settings.GetSettings()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	chatApi.RegisterChatV1Server(s, &server{})

	log.Printf("started gRPC server at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
