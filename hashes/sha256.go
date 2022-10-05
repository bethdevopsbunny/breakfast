package hashes

import (
	"crypto/sha256"
	"encoding/hex"
)

func SHA256(text string) string {

	hash := sha256.New()
	hash.Write([]byte(text))
	hashstring := hash.Sum(nil)
	return hex.EncodeToString(hashstring[:])
}
