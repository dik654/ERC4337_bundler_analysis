package main

import (
	"context"
	"log"
	"time"

	pb "github.com/dik654/Go_projects/gRPC/qrcode_microservice/go_server/proto"
)

func callGeneratePrivateKey(client pb.OtpAuthenticatorClient, id string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GeneratePrivateKey(ctx, &pb.GeneratePrivateKeyRequest{
		Id: id,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("%s", res.PrivateKey)
}

func callGenerateOtp(client pb.OtpAuthenticatorClient, id string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GenerateOtp(ctx, &pb.GenerateOtpRequest{
		Id: id,
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("%s", res.Otp)
	return res.Otp
}
