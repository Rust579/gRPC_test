package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "gRPC_test/proto/grpc_test/proto"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
	age         = 2024
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name, Age: age})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	r2, err := c.SayHi(ctx, &pb.HiRequest{Name: name, Age: age})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())
	log.Printf("Greeting: %s", r2.GetMessage())
}
