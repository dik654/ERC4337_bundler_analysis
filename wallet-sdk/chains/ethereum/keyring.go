package keyring

import (
	"crypto/ecdsa"
	"crypto/rand"

	"github.com/btcsuite/btcd/btcec/v2"
)

// 개인키 생성
func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	// btcec는 secp256k1 곡선 인스턴스 생성
	// rand.Reader는 랜덤값
	return ecdsa.GenerateKey(btcec.S256(), rand.Reader)
}
