package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	pb "gRPC_test/proto/grpc_test/proto"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50000"
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
	ctx, cancel := context.WithCancel(context.Background())
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

	stream, err := c.Chat(ctx)
	if err != nil {
		log.Fatalf("could not start chat: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a message: %v", err)
			}
			log.Printf("Received message from server: %s", in.Message)
		}
	}()

	for _, message := range []string{"Hello", "How are you?", "Goodbye"} {
		if err := stream.Send(&pb.ChatMessage{User: defaultName, Message: message}); err != nil {
			log.Fatalf("Failed to send a message: %v", err)
		}
		time.Sleep(2 * time.Second)
	}
	stream.CloseSend()
	<-waitc
}
