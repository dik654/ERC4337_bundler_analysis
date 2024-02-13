package main

import (
	"context"
	"log"
	"time"

	pb "github.com/dik654/Go_projects/gRPC/qrcode_microservice/go_server/pb"
)

func callGeneratePrivateKey(client pb.OtpAuthenticatorClient, id string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GeneratePrivateKey(ctx, &pb.GeneratePrivateKeyRequest{
		Id: id,
	})
	if err != nil {
		log.Fatalf("could not callGeneratePrivateKey: %v", err)
	}
	log.Printf("%s", res.Url)
}

func callGenerateOtp(client pb.OtpAuthenticatorClient, id string) string {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GenerateOtp(ctx, &pb.GenerateOtpRequest{
		Id: id,
	})
	if err != nil {
		log.Fatalf("could not GenerateOtp: %v", err)
	}
	log.Printf("%s", res.Otp)
	return res.Otp
}

func callVerifyOtp(client pb.OtpAuthenticatorClient, id string, otp string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.VerifyOtp(ctx, &pb.GenerateVerifyOtpRequest{
		Id:  id,
		Otp: otp,
	})
	if err != nil {
		log.Fatalf("could not callVerifyOtp: %v", err)
	}
	log.Printf("%t", res.Verification)
	return res.Verification
}
