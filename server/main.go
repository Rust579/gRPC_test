package main

import (
	"context"
	"log"
	"net"
	"strconv"

	pb "gRPC_test/proto/grpc_test/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name + " " + strconv.Itoa(int(in.Age))}, nil
}

func (s *server) SayHi(ctx context.Context, in *pb.HiRequest) (*pb.HiReply, error) {
	return &pb.HiReply{Message: "Hi " + in.Name + " " + strconv.Itoa(int(in.Age))}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
