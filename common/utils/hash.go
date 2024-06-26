package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashSHA256(name string) string {
	hasher := sha256.New()
	hasher.Write([]byte(name))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)
}
