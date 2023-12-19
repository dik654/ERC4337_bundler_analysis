package main

import (
	"context"
)

type OtpAuthenticator interface {
	GeneratePrivateKey(context.Context, string) (string, error)
	GenerateOtp(context.Context, string) (string, error)
}

type otpAuthenticator struct{}

func (s *otpAuthenticator) GeneratePrivateKey(ctx context.Context, id string) (string, error) {
	return GeneratePrivateKey(ctx, id)
}

func (s *otpAuthenticator) GenerateOtp(ctx context.Context, id string) (string, error) {
	return GenerateOtp(ctx, id)
}

func GeneratePrivateKey(ctx context.Context, session string) (string, error) {
	return "", nil
}
