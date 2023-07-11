package main

import (
	"context"
	"grpc-go/pb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:9001", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer connection.Close()

	client := pb.NewHelloServiceClient(connection)

	request := &pb.HelloRequest{Name: "Seu nome"}

	response, err := client.Hello(context.Background(), request)

	if err != nil {
		log.Fatalf("Failed to call: %v", err)
	}

	log.Printf("Response: %v", response.Msg)
}
