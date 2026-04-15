package utils

import (
	"crypto/rand"
)

func GenerateKey() {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
	}
	
}