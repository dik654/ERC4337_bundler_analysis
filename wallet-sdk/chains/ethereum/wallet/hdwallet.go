package wallet

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/cosmos/go-bip39"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// mnemonic/purpose'/ coin_type'/account'/change/address_index
// defaultHDpath := "m/44'/60'/0'/0/0"
var DefaultRootDerivationPath = accounts.DefaultRootDerivationPath
var DefaultBaseDerivationPath = accounts.DefaultBaseDerivationPath

// 이 값이 없으면 btcd/chaincfg에서 masterKey를 제대로 받아올 수 없다
const issue179FixEnvar = "GO_ETHEREUM_HDWALLET_FIX_ISSUE_179"

func (w *Wallet) SetFixIssue172(fixIssue172 bool) {
	w.fixIssue172 = fixIssue172
}

type Wallet struct {
	mnemonic    string
	masterKey   *hdkeychain.ExtendedKey
	seed        []byte
	url         accounts.URL
	paths       map[common.Address]accounts.DerivationPath
	accounts    []accounts.Account
	stateLock   sync.RWMutex
	fixIssue172 bool
}

func GenerateMnemonic() string {
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
	return mnemonic
}

// 입력받은 니모닉을 이용해서 wallet 객체 생성
func NewFromMnemonic(mnemonic string, passOpt ...string) (*Wallet, error) {
	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}

	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, errors.New("mnemonic is invalid")
	}

	seed, err := NewSeedFromMnemonic(mnemonic, passOpt...)
	if err != nil {
		return nil, err
	}

	wallet, err := newWallet(seed)
	if err != nil {
		return nil, err
	}
	wallet.mnemonic = mnemonic

	return wallet, nil
}

// 니모닉에서 시드 생성
func NewSeedFromMnemonic(mnemonic string, passOpt ...string) ([]byte, error) {
	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}

	password := ""
	if len(passOpt) > 0 {
		password = passOpt[0]
	}

	return bip39.NewSeedWithErrorChecking(mnemonic, password)
}

// 체인에 맞게 masterKey를 받아와서 wallet 구조체에 저장
func newWallet(seed []byte) (*Wallet, error) {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		masterKey:   masterKey,
		seed:        seed,
		accounts:    []accounts.Account{},
		paths:       map[common.Address]accounts.DerivationPath{},
		fixIssue172: false || len(os.Getenv(issue179FixEnvar)) > 0,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////

// derive path 검사 에러 체크를 강제화시키는 함수
func MustParseDerivationPath(path string) accounts.DerivationPath {
	parsed, err := accounts.ParseDerivationPath(path)
	if err != nil {
		panic(err)
	}

	return parsed
}

// 인수로 들어온 derive path가 정상적인 값인지 체크
// 정상적이라면 공백문자 등을 제거한 뒤 리턴
func ParseDerivationPath(path string) (accounts.DerivationPath, error) {
	return accounts.ParseDerivationPath(path)
}

// 인수로 들어온 derive path로부터 개인키 -> 공개키 -> 지갑주소 유도
// pin이 하는 역할은 잘 모르겠음
func (w *Wallet) Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error) {
	// Try to derive the actual account and update its URL if successful
	w.stateLock.RLock() // Avoid device disappearing during derivation

	address, err := w.deriveAddress(path)

	w.stateLock.RUnlock()

	// If an error occurred or no pinning was requested, return
	if err != nil {
		return accounts.Account{}, err
	}

	account := accounts.Account{
		Address: address,
		URL: accounts.URL{
			Scheme: "",
			Path:   path.String(),
		},
	}

	if !pin {
		return account, nil
	}

	// Pinning needs to modify the state
	w.stateLock.Lock()
	defer w.stateLock.Unlock()

	if _, ok := w.paths[address]; !ok {
		w.accounts = append(w.accounts, account)
		w.paths[address] = path
	}

	return account, nil
}

// derive path로부터 개인키 유도
func (w *Wallet) derivePrivateKey(path accounts.DerivationPath) (*ecdsa.PrivateKey, error) {
	var err error
	key := w.masterKey
	// 이 부분이 없으면 masterKey가 제대로 안나옴
	for _, n := range path {
		if w.fixIssue172 && key.IsAffectedByIssue172() {
			key, err = key.Derive(n)
		} else {
			key, err = key.DeriveNonStandard(n)
		}
		if err != nil {
			return nil, err
		}
	}

	privateKey, err := key.ECPrivKey()
	privateKeyECDSA := privateKey.ToECDSA()
	if err != nil {
		return nil, err
	}

	return privateKeyECDSA, nil
}

// derive path로부터 유도된 개인키에서 공개키 계산
func (w *Wallet) derivePublicKey(path accounts.DerivationPath) (*ecdsa.PublicKey, error) {
	privateKeyECDSA, err := w.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}

	return publicKeyECDSA, nil
}

// derive path로부터 유도된 개인키에서 계산한 공개키에서 지갑주소 계산
func (w *Wallet) deriveAddress(path accounts.DerivationPath) (common.Address, error) {
	publicKeyECDSA, err := w.derivePublicKey(path)
	if err != nil {
		return common.Address{}, err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

// account 슬라이스에서 index에 해당하는 account 삭제하는 기능
// (여러 타입의 derive path로부터 account들을 생성할 수 있으므로)
func removeAtIndex(accts []accounts.Account, index int) []accounts.Account {
	return append(accts[:index], accts[index+1:]...)
}

////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////

// wallet 인터페이스에 포함되어있는 필수 메서드들

// 하드웨어 월렛용 메서드
func (w *Wallet) Status() (string, error) {
	return "ok", nil
}

// 하드웨어 월렛용 메서드
func (w *Wallet) Open(passphrase string) error {
	return nil
}

// 하드웨어 월렛용 메서드
func (w *Wallet) Close() error {
	return nil
}

// 이더리움 쪽 라이브러리로부터 masterKey를 받아오는 메서드인데 구현되어있지 않음
func (w *Wallet) SelfDerive(base []accounts.DerivationPath, chain ethereum.ChainStateReader) {
	// TODO: self derivation
}

// pin의 용도는 잘 모르겠음
func (w *Wallet) Unpin(account accounts.Account) error {
	w.stateLock.RLock()
	defer w.stateLock.RUnlock()

	for i, acct := range w.accounts {
		if acct.Address.String() == account.Address.String() {
			w.accounts = removeAtIndex(w.accounts, i)
			delete(w.paths, account.Address)
			return nil
		}
	}

	return errors.New("account not found")
}

// derive path Getter 메서드
func (w *Wallet) URL() accounts.URL {
	return w.url
}

// HDwallet account들 Getter 메서드
// 중간에 수정되거나 하는 보안 문제를 방지하기 위해 락을 걸고
// 포인터를 리턴하는 것이 아닌 값을 복사해서 리턴
func (w *Wallet) Accounts() []accounts.Account {
	// Attempt self-derivation if it's running
	// Return whatever account list we ended up with
	w.stateLock.RLock()
	defer w.stateLock.RUnlock()

	cpy := make([]accounts.Account, len(w.accounts))
	copy(cpy, w.accounts)
	return cpy
}

// 인수로 들어온 주소가 이 wallet안에 들어있는지 체크
func (w *Wallet) Contains(account accounts.Account) bool {
	w.stateLock.RLock()
	defer w.stateLock.RUnlock()

	_, exists := w.paths[account.Address]
	return exists
}

//
// 데이터 서명 메서드
//

// 인수로 들어온 byte 슬라이스 데이터를 keccak256으로 해싱한 뒤 서명하는 메서드
func (w *Wallet) SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error) {
	// Make sure the requested account is contained within
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}

	return w.SignHash(account, crypto.Keccak256(data))
}

// passphrase를 이용하여 keccak256 해시에 사인하는 것은 불가능한 걸로 보임
func (w *Wallet) SignDataWithPassphrase(account accounts.Account, passphrase, mimeType string, data []byte) ([]byte, error) {
	// Make sure the requested account is contained within
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}

	return w.SignHashWithPassphrase(account, passphrase, crypto.Keccak256(data))
}

// 해시가 아니더라도 인수로 들어온 임의의 byte슬라이스를 서명하는 메서드
func (w *Wallet) SignText(account accounts.Account, text []byte) ([]byte, error) {
	// Make sure the requested account is contained within
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}

	return w.SignHash(account, accounts.TextHash(text))
}

func (w *Wallet) SignTextWithPassphrase(account accounts.Account, passphrase string, text []byte) ([]byte, error) {
	// Make sure the requested account is contained within
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}

	return w.SignHashWithPassphrase(account, passphrase, accounts.TextHash(text))
}

// passphrase를 이용하여 해시를 사인하는 것은 불가능한 걸로 보임
func (w *Wallet) SignHashWithPassphrase(account accounts.Account, passphrase string, hash []byte) ([]byte, error) {
	return w.SignHash(account, hash)
}

// 인수로 들어온 해시값 account로 서명하는 메서드
func (w *Wallet) SignHash(account accounts.Account, hash []byte) ([]byte, error) {
	// Make sure the requested account is contained within
	path, ok := w.paths[account.Address]
	if !ok {
		return nil, accounts.ErrUnknownAccount
	}

	privateKey, err := w.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}

	return crypto.Sign(hash, privateKey)
}

//
// 트랜잭션 서명 메서드
//

// passphrase를 넣고 트랜잭션을 서명할 수는 없는 걸로 보임
func (w *Wallet) SignTxWithPassphrase(account accounts.Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return w.SignTx(account, tx, chainID)
}

// account를 이용하여 트랜잭션을 서명하는 메서드
func (w *Wallet) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	w.stateLock.RLock() // Comms have own mutex, this is for the state fields
	defer w.stateLock.RUnlock()

	// Make sure the requested account is contained within
	path, ok := w.paths[account.Address]
	if !ok {
		return nil, accounts.ErrUnknownAccount
	}

	privateKey, err := w.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}

	signer := types.LatestSignerForChainID(chainID)

	// Sign the transaction and verify the sender to avoid hardware fault surprises
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return nil, err
	}

	sender, err := types.Sender(signer, signedTx)
	if err != nil {
		return nil, err
	}

	if sender != account.Address {
		return nil, fmt.Errorf("signer mismatch: expected %s, got %s", account.Address.Hex(), sender.Hex())
	}

	return signedTx, nil
}

//
// 아래는 account로부터 개인키를 가져오는 메서드들
//

// 계산된 binary 개인키를 hex값으로 변환
func (w *Wallet) PrivateKeyHex(account accounts.Account) (string, error) {
	privateKeyBytes, err := w.PrivateKeyBytes(account)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(privateKeyBytes)[2:], nil
}

// 유도된 개인키 구조체로부터 binary 개인키의 계산
func (w *Wallet) PrivateKeyBytes(account accounts.Account) ([]byte, error) {
	privateKey, err := w.PrivateKey(account)
	if err != nil {
		return nil, err
	}

	return crypto.FromECDSA(privateKey), nil
}

// account에 있는 derive path(URL)에서 개인키 구조체 유도
func (w *Wallet) PrivateKey(account accounts.Account) (*ecdsa.PrivateKey, error) {
	path, err := ParseDerivationPath(account.URL.Path)
	if err != nil {
		return nil, err
	}

	return w.derivePrivateKey(path)
}

//////
///////

// SignTxEIP155 implements accounts.Wallet, which allows the account to sign an ERC-20 transaction.
func (w *Wallet) SignTxEIP155(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	w.stateLock.RLock() // Comms have own mutex, this is for the state fields
	defer w.stateLock.RUnlock()

	// Make sure the requested account is contained within
	path, ok := w.paths[account.Address]
	if !ok {
		return nil, accounts.ErrUnknownAccount
	}

	privateKey, err := w.derivePrivateKey(path)
	if err != nil {
		return nil, err
	}

	signer := types.NewEIP155Signer(chainID)
	// Sign the transaction and verify the sender to avoid hardware fault surprises
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return nil, err
	}

	sender, err := types.Sender(signer, signedTx)
	if err != nil {
		return nil, err
	}

	if sender != account.Address {
		return nil, fmt.Errorf("signer mismatch: expected %s, got %s", account.Address.Hex(), sender.Hex())
	}

	return signedTx, nil
}
