package main

import (
	"log"

	pb "github.com/dik654/Go_projects/gRPC/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":8080"
)

func main() {
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreetServiceClient(conn)

	names := &pb.NamesList{
		Names: []string{"gildong", "nick", "alia"},
	}

	// callSayHello(client)
	// callSayHelloServerStream(client, names)
	callSayHelloClientStream(client, names)
	// callHelloBidirectionalStream(client, names)
}
