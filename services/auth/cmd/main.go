package main

import (
  "context"
  "fmt"
  "log"
  "net"

  gofakeit "github.com/brianvoe/gofakeit/v7"
  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"
  "google.golang.org/protobuf/types/known/emptypb"
  "google.golang.org/protobuf/types/known/timestamppb"

  // use general reference for easier api version change in code
  userApi "github.com/Genvekt/cli-chat/libraries/api/user/v1"
  "github.com/Genvekt/cli-chat/services/auth/settings"
)

type server struct {
  userApi.UnimplementedUserV1Server
}

func (s *server) Get(_ context.Context, req *userApi.GetRequest) (*userApi.GetResponse, error) {
  log.Printf("Received get user request: %+v", req)

  return &userApi.GetResponse{
    Id:        req.Id,
    Name:      gofakeit.Name(),
    Email:     gofakeit.Email(),
    Role:      userApi.UserRole_USER,
    CreatedAt: timestamppb.New(gofakeit.Date()),
    UpdatedAt: timestamppb.New(gofakeit.Date()),
  }, nil
}

func (s *server) Create(_ context.Context, req *userApi.CreateRequest) (*userApi.CreateResponse, error) {
  log.Printf("Received create user request: %+v", req)
  return &userApi.CreateResponse{
    Id: gofakeit.Int64(),
  }, nil
}

func (s *server) Update(_ context.Context, req *userApi.UpdateRequest) (*emptypb.Empty, error) {
  log.Printf("Received update user request: %+v", req)
  return &emptypb.Empty{}, nil
}

func (s *server) Delete(_ context.Context, req *userApi.DeleteRequest) (*emptypb.Empty, error) {
  log.Printf("Received delete user request: %+v", req)
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
  userApi.RegisterUserV1Server(s, &server{})

  log.Printf("started gRPC server at %v", lis.Addr())
  if err := s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}
