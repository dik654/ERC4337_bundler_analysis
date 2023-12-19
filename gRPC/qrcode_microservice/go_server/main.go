package main

import (
	"log"
	"net"

	pb "github.com/dik654/Go_projects/gRPC/qrcode_microservice/go_server/proto"
	"google.golang.org/grpc"
)

type OtpAuthenticatorServer struct {
	pb.OtpAuthenticatorServer
	authenticator OtpAuthenticator
}

func NewOtpAuthenticatorServer(authenticator OtpAuthenticator) *OtpAuthenticatorServer {
	return &OtpAuthenticatorServer{
		authenticator: authenticator,
	}
}

func main() {
	secrets := make(map[string]string)
	authenticator := NewLoggingService(NewMetricService(&otpAuthenticator{secrets}))

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("Failed to start the server %v", err)
	}

	grpcServer := grpc.NewServer()
	otpAuthenticatorServer := NewOtpAuthenticatorServer(authenticator)
	pb.RegisterGreetServiceServer(grpcServer, otpAuthenticatorServer)
	log.Printf("server started at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

}
