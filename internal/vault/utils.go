package vault

import "crypto/sha256"

func sha(b []byte) []byte {
	h := sha256.Sum256(b)
	return h[:]
}
