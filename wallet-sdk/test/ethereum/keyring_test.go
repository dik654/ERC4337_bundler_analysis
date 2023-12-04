package test

import (
	"crypto/ecdsa"
	"testing"

	keyring "github.com/dik654/Go_projects/wallet-sdk/chains/ethereum"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	// 개인 키 생성
	privateKey, err := keyring.GeneratePrivateKey()

	// 개인키 생성 도중 에러가 없는지 체크
	assert.NoError(t, err)

	// 개인 키가 nil인 경우가 있는지 체크
	assert.NotNil(t, privateKey)

	// 공개 키가 유효한지 확인
	publicKey := privateKey.Public()
	assert.NotNil(t, publicKey)

	// 공개 키 타입이 *ecdsa.PublicKey인지 확인
	_, ok := publicKey.(*ecdsa.PublicKey)
	assert.True(t, ok)
}
