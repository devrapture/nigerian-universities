package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
)

func GenerateRawKey() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return "sk_live_" + hex.EncodeToString(bytes), nil
}

func HashKey(rawKey string) string {
	h := sha256.Sum256([]byte(rawKey))

	return hex.EncodeToString(h[:])
}

func VerifyKey(rawKey, hashedKey string) bool {
	computedHash := HashKey(rawKey)

	return subtle.ConstantTimeCompare([]byte(computedHash), []byte(hashedKey)) == 1
}
