package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func CombineSessionDataAndSignature(sessionData string, signature string) []byte {
	separator := []byte("|")
	combined := append([]byte(sessionData), separator...)
	combined = append(combined, []byte(signature)...)
	return combined
}

func CreateSignature(value string, secretKey string) string {
	hmac := hmac.New(sha256.New, []byte(secretKey))
	hmac.Write([]byte(value))
	return hex.EncodeToString(hmac.Sum(nil))
}
