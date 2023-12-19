package main

import (
	"log"
	"net"

	pb "github.com/dik654/Go_projects/gRPC/qrcode_microservice/go_server/pb"
	"google.golang.org/grpc"
)

func main() {
	secrets := make(map[string]string)
	authenticator := NewLoggingService(NewMetricService(&otpAuthenticator{secrets}))

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("Failed to start the server %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOtpAuthenticatorServer(grpcServer, NewOtpAuthenticatorServer(authenticator))
	log.Printf("server started at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

}
