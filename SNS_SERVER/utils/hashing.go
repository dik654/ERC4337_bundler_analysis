package utils

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"math/big"
	"os"
	"strings"

	"github.com/xdg-go/pbkdf2"
)

const (
	IterationCount = 25000
)

func generateRandomSalt(length int) ([]byte, error) {
	results := make([]byte, length)
	for i := 0; i < length; i++ {
		salt, err := rand.Int(rand.Reader, big.NewInt(255))
		if err != nil {
			return nil, err
		}
		results[i] = byte(salt.Int64())
	}
	return results, nil
}

func HashingPassword(password string) string {
	salt := []byte(os.Getenv("PDKDF2_SALT"))
	iteration := IterationCount
	hashedPassword := GeneratePBKDF2([]byte(password), salt, iteration, sha512.Size)
	return string(hashedPassword)
}

func GeneratePBKDF2(password []byte, salt []byte, iterateCount int, keyLength int) string {
	key := pbkdf2.Key(password, salt, iterateCount, keyLength, sha512.New)
	hexValue := strings.ToUpper(hex.EncodeToString(key))
	return hexValue
}
