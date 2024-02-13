package main

import (
	"context"
	"fmt"
	"time"

	"github.com/xlzd/gotp"
)

type OtpAuthenticator interface {
	GeneratePrivateKey(context.Context, string) (string, error)
	GenerateOtp(context.Context, string) (string, error)
	VerifyOtp(context.Context, string, string) (bool, error)
}

type otpAuthenticator struct {
	secrets map[string]string
}

func (s *otpAuthenticator) GeneratePrivateKey(ctx context.Context, id string) (string, error) {
	uri, randomSecret, err := GeneratePrivateKey(ctx, id)
	if err != nil {
		return "", err
	}
	s.secrets[id] = randomSecret
	return uri, nil
}

func (s *otpAuthenticator) GenerateOtp(ctx context.Context, id string) (string, error) {
	secret, exist := s.secrets[id]
	if !exist {
		return "", fmt.Errorf("no secret key found for ID: %s", id)
	}
	return GenerateOtp(ctx, secret)
}

func (s *otpAuthenticator) VerifyOtp(ctx context.Context, id string, otp string) (bool, error) {
	secret, exist := s.secrets[id]
	if !exist {
		return false, fmt.Errorf("no secret key found for ID: %s", id)
	}
	return VerifyOtp(ctx, secret, otp), nil
}

func GeneratePrivateKey(ctx context.Context, id string) (string, string, error) {
	randomSecret := gotp.RandomSecret(16)
	uri := fmt.Sprintf("otpauth://totp/authentication:%s?secret=%s&issuer=DIK", id, randomSecret)
	return uri, randomSecret, nil
}

func GenerateOtp(ctx context.Context, secret string) (string, error) {
	totp := gotp.NewDefaultTOTP(secret)
	return totp.Now(), nil
}

func VerifyOtp(ctx context.Context, secret string, otp string) bool {
	totp := gotp.NewDefaultTOTP(secret)

	ok := totp.Verify(otp, time.Now().Unix())
	return ok
}
