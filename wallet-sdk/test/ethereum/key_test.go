package test

import (
	"crypto/ecdsa"
	"testing"

	key "github.com/dik654/Go_projects/wallet-sdk/chains/ethereum/wallet"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	// 개인 키 생성
	privateKey, err := key.GeneratePrivateKey()

	// 개인키 생성 도중 에러가 없는지 체크
	assert.NoError(t, err)

	// 개인 키가 nil인 경우가 있는지 체크
	assert.NotNil(t, privateKey)

	// 공개 키가 유효한지 체크
	publicKey := privateKey.Public()
	assert.NotNil(t, publicKey)

	// 공개 키 타입이 *ecdsa.PublicKey을 만족하는지 체크
	_, ok := publicKey.(*ecdsa.PublicKey)
	assert.True(t, ok)
}
