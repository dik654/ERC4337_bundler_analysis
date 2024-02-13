package main

import (
	"context"
	"fmt"

	pb "github.com/dik654/Go_projects/gRPC/qrcode_microservice/go_server/pb"
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

func (s *OtpAuthenticatorServer) GeneratePrivateKey(ctx context.Context, req *pb.GeneratePrivateKeyRequest) (*pb.GeneratePrivateKeyResponse, error) {
	url, err := s.authenticator.GeneratePrivateKey(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GeneratePrivateKeyResponse{Url: url}, nil
}

func (s *OtpAuthenticatorServer) GenerateOtp(ctx context.Context, req *pb.GenerateOtpRequest) (*pb.GenerateOtpRequestResponse, error) {
	otp, err := s.authenticator.GenerateOtp(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GenerateOtpRequestResponse{Otp: otp}, nil
}

func (s *OtpAuthenticatorServer) VerifyOtp(ctx context.Context, req *pb.GenerateVerifyOtpRequest) (*pb.GenerateVerifyOtpResponse, error) {
	fmt.Println(req)
	ok, err := s.authenticator.VerifyOtp(ctx, req.Id, req.Otp)
	if err != nil {
		return nil, err
	}
	return &pb.GenerateVerifyOtpResponse{Verification: ok}, nil
}
