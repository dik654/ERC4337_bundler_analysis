package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtGenerator interface {
	GenerateJWT(context.Context, string) error
}

type jwtGenerator struct{}

func (s *jwtGenerator) GenerateJWT(ctx context.Context, userId string, session string) (string, error) {
	return GenerateJWT(ctx, userId, session)
}

func GenerateJWT(ctx context.Context, userId string, session string) (string, error) {
	ecdsaPrivateKey, err := loadPrivateKey(os.Getenv("ECDSA_PRIVATE_KEY_PATH"))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims{
		"version": os.Getenv("QR_GENERATOR_SERVICE_VERSION"),
		"expire":  time.Now().Add(30 * time.Second).Unix(),
		"user_id": userId,
		"issuer":  "DIK",
	})

	tokenString, err := token.SignedString(ecdsaPrivateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func loadPrivateKey(relativePath string) (*ecdsa.PrivateKey, error) {
	// 현재 디렉토리의 절대 경로
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// 절대 경로와 상대 경로 합치기
	keyPath := filepath.Join(currentDir, relativePath)

	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
