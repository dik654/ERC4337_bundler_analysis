package main

import (
	"context"

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
	privateKey, err := s.authenticator.GeneratePrivateKey(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GeneratePrivateKeyResponse{PrivateKey: privateKey}, nil
}

func (s *OtpAuthenticatorServer) GenerateOtp(ctx context.Context, req *pb.GenerateOtpRequest) (*pb.GenerateOtpRequestResponse, error) {
	otp, err := s.authenticator.GenerateOtp(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GenerateOtpRequestResponse{Otp: otp}, nil
}
