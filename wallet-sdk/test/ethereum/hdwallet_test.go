package test

import (
	"testing"

	"github.com/dik654/Go_projects/wallet-sdk/chains/ethereum/wallet"
	"github.com/ethereum/go-ethereum/accounts"
)

func TestWallet(t *testing.T) {
	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	myWallet, err := wallet.NewFromMnemonic(mnemonic)
	if err != nil {
		t.Error(err)
	}
	// Check that Wallet implements the accounts.Wallet interface.
	var _ accounts.Wallet = myWallet

	path, err := wallet.ParseDerivationPath("m/44'/60'/0'/0/0")
	if err != nil {
		t.Error(err)
	}

	account, err := myWallet.Derive(path, false)
	if err != nil {
		t.Error(err)
	}

	if account.Address.Hex() != "0xC49926C4124cEe1cbA0Ea94Ea31a6c12318df947" {
		t.Error("wrong address")
	}

	if len(myWallet.Accounts()) != 0 {
		t.Error("expected 0")
	}

	account, err = myWallet.Derive(path, true)
	if err != nil {
		t.Error(err)
	}

	if len(myWallet.Accounts()) != 1 {
		t.Error("expected 1")
	}
}
