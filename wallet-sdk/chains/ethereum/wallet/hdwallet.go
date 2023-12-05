package wallet

import (
	"fmt"
	"log"

	"github.com/cosmos/go-bip39"
)

// mnemonic/purpose'/ coin_type'/account'/change/address_index

func GenerateMnemonic() string {
	// defaultHDpath := "m/44'/60'/0'/0/0"
	// 128bits 랜덤값 생성
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		log.Fatalf("Error generating entropy: %v", err)
	}
	// 랜덤값 기반으로 니모닉 생성 (랜덤값 + 해시된 MAC)
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		log.Fatalf("Error generating mnemonic: %v", err)
	}
	fmt.Printf(mnemonic)
	return mnemonic
}
