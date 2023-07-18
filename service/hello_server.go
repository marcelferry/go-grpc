package service

import (
	"context"

	"grpc-go/pb"
)

type HelloServer struct {
	pb.UnimplementedHelloServiceServer
}

// NewHelloServer returns a new hello server
func NewHelloServer() pb.HelloServiceServer {
	return &HelloServer{}
}

func (s *HelloServer) Hello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Msg: "Hello " + request.GetName()}, nil
}

func (s *HelloServer) Admin(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Msg: "Admin Hello " + request.GetName()}, nil
}

func (s *HelloServer) Protected(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Msg: "Protected Hello " + request.GetName()}, nil
}
