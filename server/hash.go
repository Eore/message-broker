package server

import (
	"crypto/sha512"
	"encoding/hex"
)

func Hashing(s string) string {
	h := sha512.New512_256()
	// h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
