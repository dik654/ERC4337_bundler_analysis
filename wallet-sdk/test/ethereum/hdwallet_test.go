package test

import (
	"testing"

	hdwallet "github.com/dik654/Go_projects/wallet-sdk/chains/ethereum/wallet"
)

func TestGenerateMnemonic(t testing.T) {
	hdwallet.GenerateMnemonic()
}
