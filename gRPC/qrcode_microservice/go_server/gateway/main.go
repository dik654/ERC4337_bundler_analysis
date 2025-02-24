package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"

	gw "github.com/dik654/Go_projects/gRPC/qrcode_microservice/go_server/pb"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:3000", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	conn, err := grpc.DialContext(ctx, *grpcServerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	mux := runtime.NewServeMux()
	err = gw.RegisterOtpAuthenticatorHandler(ctx, mux, conn)
	if err != nil {
		return err
	}
	httpListener, err := net.Listen("tcp", ":3002")
	if err != nil {
		return fmt.Errorf("Failed to listen on %v: %v", httpListener.Addr(), err)
	}
	log.Printf("Gateway server started at %v", httpListener.Addr())
	return http.Serve(httpListener, mux)
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		grpclog.Fatal(err)
	}
}
