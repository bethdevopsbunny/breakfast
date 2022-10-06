package hashes

import (
	"crypto/sha1"
	"encoding/hex"
)

func SHA1(text string) string {

	hash := sha1.New()
	hash.Write([]byte(text))
	hashstring := hash.Sum(nil)
	return hex.EncodeToString(hashstring[:])
}
