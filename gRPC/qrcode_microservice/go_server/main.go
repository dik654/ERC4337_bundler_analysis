package main

import (
	"log"
	"net"
	"net/http"

	pb "github.com/dik654/Go_projects/gRPC/qrcode_microservice/go_server/pb"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {

	secrets := make(map[string]string)
	authenticator := NewLoggingService(NewMetricService(&otpAuthenticator{secrets}))

	go func() {
		httpListener, err := net.Listen("tcp", ":3001")
		if err != nil {
			log.Fatalf("Failed to listen on port 3001: %v", err)
		}
		log.Printf("Prometheus metrics server started at %v", httpListener.Addr())
		http.Handle("/metrics", promhttp.Handler())
		if err := http.Serve(httpListener, nil); err != nil {
			log.Fatalf("Failed to start the prometheus metrics server: %v", err)
		}
	}()

	grpcListener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("Failed to start the server %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterOtpAuthenticatorServer(grpcServer, NewOtpAuthenticatorServer(authenticator))
	log.Printf("gRPC server started at %v", grpcListener.Addr())
	if err := grpcServer.Serve(grpcListener); err != nil {
		log.Fatalf("Failed to start: %v", err)
	}
}
