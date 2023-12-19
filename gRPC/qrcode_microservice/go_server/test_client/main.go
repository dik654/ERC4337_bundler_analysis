package main

import (
	"log"

	pb "github.com/dik654/Go_projects/gRPC/qrcode_microservice/go_server/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	port = ":3000"
)

func main() {
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOtpAuthenticatorClient(conn)

	callGeneratePrivateKey(client, "test")
	callGenerateOtp(client, "test")
}
