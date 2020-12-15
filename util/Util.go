package util

import (
	"crypto/sha1"
	"fmt"
)

func HashSha1Hex(input []byte) string {
	h := sha1.New()
	h.Write(input)
	result := fmt.Sprintf("%x", h.Sum(nil))
	return result
}
