package main

import (
	"context"
	"fmt"

	"github.com/xlzd/gotp"
)

type OtpAuthenticator interface {
	GeneratePrivateKey(context.Context, string) (string, error)
	GenerateOtp(context.Context, string) (string, error)
}

type otpAuthenticator struct {
	secrets map[string]string
}

func (s *otpAuthenticator) GeneratePrivateKey(ctx context.Context, id string) (string, error) {
	randomSecret, err := GeneratePrivateKey(ctx, id)
	if err != nil {
		return "", err
	}
	s.secrets[id] = randomSecret
	return randomSecret, nil
}

func (s *otpAuthenticator) GenerateOtp(ctx context.Context, id string) (string, error) {
	secret, exist := s.secrets[id]
	if !exist {
		return "", fmt.Errorf("no secret key found for ID: %s", id)
	}
	return GenerateOtp(ctx, secret)
}

func GeneratePrivateKey(ctx context.Context, id string) (string, error) {
	randomSecret := gotp.RandomSecret(16)
	return randomSecret, nil
}

func GenerateOtp(ctx context.Context, secret string) (string, error) {
	totp := gotp.NewDefaultTOTP(secret)
	return totp.Now(), nil
}
